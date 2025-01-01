// go:build ignore
// +build ignore

// The printk function is used to print messages to the kernel log. It is
// similar to the printf function in C. `sudo cat
// /sys/kernel/debug/tracing/trace_pipe`

#include <asm-generic/int-ll64.h>
#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <stdlib.h>

#define MAX_PAYLOAD_SIZE 2048

struct packet_s {
  long payload_size;
  unsigned char payload[MAX_PAYLOAD_SIZE];
};

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 4096);
} module_map SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_ARRAY);
  __uint(max_entries, 1);
  __type(key, __u32);
  __type(value, _Bool);
} control_map SEC(".maps");

SEC("xdp")
int share_data(struct xdp_md *ctx) {

  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;

  if (data + sizeof(struct ethhdr) > data_end) {
    return XDP_PASS;
  }

  // Calcola la dimensione del payload
  long payload_size = data_end - data;
  if (payload_size > MAX_PAYLOAD_SIZE) {
    bpf_printk("Payload size too big: %ld\n", payload_size);
    payload_size = MAX_PAYLOAD_SIZE;
  }

  struct packet_s *packet =
      bpf_ringbuf_reserve(&module_map, sizeof(struct packet_s), 0);
  if (!packet) {
    return XDP_PASS;
  }

  packet->payload_size = payload_size;

  unsigned char *src = data;
  for (long i = 0; i < payload_size; i++) {
    if (src + i >= data_end) {
      bpf_ringbuf_discard(packet, 0);
      return XDP_ABORTED;
    }
    packet->payload[i] = src[i];
  }

  bpf_ringbuf_submit(packet, 0);

  bpf_printk("Packet processed\n");

  __u32 key = 0;
  _Bool *control_value = bpf_map_lookup_elem(&control_map, &key);
  if (control_value) {
    if (*control_value) {
      bpf_printk("Passing package\n");
      return XDP_PASS;
    } else {
      bpf_printk("Dropping package\n");
      return XDP_DROP;
    }
  } else {
    bpf_printk("Failed to retrieve control value from map\n");
    return XDP_PASS;
  }
}

char __license[] SEC("license") = "Dual MIT/GPL";
