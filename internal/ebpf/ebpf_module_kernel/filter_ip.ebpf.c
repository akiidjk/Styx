#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>

#define MAX_PAYLOAD_SIZE 2048
#define MAX_FILTERS 1024

struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(max_entries, 16);
  __uint(pinning, LIBBPF_PIN_BY_NAME);
} event_output SEC(".maps");

// Structs
struct log_event {
  __u32 src_ip;
  __u32 filter_ip;
  char message[128];
};

// Maps
struct {
  __uint(type, BPF_MAP_TYPE_ARRAY);
  __type(key, __u32);
  __type(value, __u32);
  __uint(max_entries, 1024);
} ip_filter_map SEC(".maps");

SEC("xdp")
int xdp_filter_ip(struct xdp_md *ctx) {
  bpf_printk("FILTERING....\n");

  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;
  __u32 key = 0;
  __u32 *value;

  struct log_event event = {};
  struct ethhdr *eth = data;

  if ((void *)(eth + 1) > data_end)
    return XDP_DROP;

  if (eth->h_proto != __constant_htons(ETH_P_IP))
    return XDP_PASS;

  struct iphdr *iph = (void *)(eth + 1);
  if ((void *)(iph + 1) > data_end)
    return XDP_DROP;

  __u32 src_ip = iph->saddr;

#pragma unroll
  for (int i = 0; i < MAX_FILTERS; i++) {
    key = i;
    value = bpf_map_lookup_elem(&ip_filter_map, &key);

    if (value == NULL) {
      break;
    }

    if (*value == 0) {
      break;
    }

    event.src_ip = src_ip;
    event.filter_ip = *value;
    __builtin_memcpy(event.message, "Valori degli IP", 15);
    bpf_perf_event_output(ctx, &event_output, BPF_F_CURRENT_CPU, &event,
                          sizeof(event));

    if (src_ip == *value) {
      return XDP_DROP;
    }
  }

  return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";
