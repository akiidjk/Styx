#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>

// struct log_event {
//   __u64 timestamp;
//   __u32 src_mac;
//   __u32 filter_mac;
//   char message[128];
// };

// // Maps
// struct {
//   __uint(type, BPF_MAP_TYPE_ARRAY);
//   __type(key, __u32);
//   __type(value, __u32);
//   __uint(max_entries, 1024);
// } mac_filter_map SEC(".maps");

// struct {
//   __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
//   __uint(max_entries, 0);
// } event_output_map SEC(".maps");

SEC("xdp")
int xdp_filter_mac(struct xdp_md *ctx) {
  bpf_printk("FILTERING....MAC \n");
  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;
  __u32 key = 0;
  __u32 *value;

  // struct log_event event = {};
  struct ethhdr *eth = data;

  if ((void *)eth + sizeof(*eth) > data_end)
    return XDP_PASS;

  unsigned char *src_mac = eth->h_source;
  unsigned char *dst_mac = eth->h_dest;

  // For debugging purposes,
  bpf_printk("Source MAC: %02x:%02x:%02x:%02x:%02x:%02x\n", src_mac[0],
             src_mac[1], src_mac[2], src_mac[3], src_mac[4], src_mac[5]);
  bpf_printk("Destination MAC: %02x:%02x:%02x:%02x:%02x:%02x\n", dst_mac[0],
             dst_mac[1], dst_mac[2], dst_mac[3], dst_mac[4], dst_mac[5]);

  return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";
