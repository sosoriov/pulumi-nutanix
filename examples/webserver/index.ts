import * as nutanix from "@pulumi/nutanix";

const vm = new nutanix.VirtualMachine("vm", {
    description: "webserver test from pulumi",
    num_sockets: 1,
    num_vcpus_per_socket: 2,
    memory_size_mib: 4096
});

export const ip = vm.ip_address

