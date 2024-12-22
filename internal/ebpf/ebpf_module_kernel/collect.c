// go:build ignore
// +build ignore

#include <bpf/bpf_endian.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>

struct packet_info {
  __u32 src_ip;
  __u32 dst_ip;
  __u16 protocol;
  __u16 pad;
};

struct {
  __uint(type, BPF_MAP_TYPE_RINGBUF);
  __uint(max_entries, 4096);
} pkt_data SEC(".maps");

SEC("xdp")
int share_data(struct xdp_md *ctx) {
  __u32 key = 0;
  struct packet_info info = {0};

  // Get the packet data
  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;

  // Some check
  struct ethhdr *eth = data;
  if ((void *)(eth + 1) > data_end)
    return XDP_PASS;

  if (eth->h_proto != bpf_htons(ETH_P_IP))
    return XDP_PASS;

  struct iphdr *ip = (void *)(eth + 1);
  if ((void *)(ip + 1) > data_end)
    return XDP_PASS;

  // Fill the packet info

  info.src_ip = ip->saddr;
  info.dst_ip = ip->daddr;
  info.protocol = ip->protocol;

  bpf_ringbuf_output(&pkt_data, &info, sizeof(info), 0);

  return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";
