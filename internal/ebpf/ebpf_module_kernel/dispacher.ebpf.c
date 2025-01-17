#include <bpf/bpf_helpers.h>
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>

struct {
  __uint(type, BPF_MAP_TYPE_PROG_ARRAY);
  __uint(max_entries, 2);
  __uint(key_size, sizeof(__u32));
  __uint(value_size, sizeof(__u32));
  __uint(pinning, LIBBPF_PIN_BY_NAME);
} prog_array SEC(".maps");

struct {
  __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
  __uint(max_entries, 16);
  __uint(pinning, LIBBPF_PIN_BY_NAME);
} event_output SEC(".maps");

SEC("xdp")
int xdp_dispach(struct xdp_md *ctx) {
  bpf_tail_call(ctx, &prog_array, 0);
  return XDP_PASS;
}
char __license[] SEC("license") = "Dual MIT/GPL";
