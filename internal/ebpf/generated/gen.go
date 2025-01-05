package ebpfGenerated

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go counter ../ebpf_module_kernel/counter.c --pkg counterebpf
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go collecter ../ebpf_module_kernel/collect.c --pkg collectebpf
