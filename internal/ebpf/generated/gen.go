package ebpfGenerated

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter ../ebpf_module_kernel/counter.c --pkg counterebpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go filterip ../ebpf_module_kernel/filter_ip.c --pkg filterip
