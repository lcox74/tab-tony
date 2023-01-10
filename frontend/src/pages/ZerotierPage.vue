<template>

    <div class="grid h-screen place-items-center">

        <Verification v-if="!Verified" :validate="OnValidate" scope="zerotier" />
        <div v-else>
            <div v-if="isLoading">
                <DotLoader />
            </div>
            <div v-else class="block rounded-lg shadow-lg bg-white max-w-lg">
                <div class="relative overflow-x-auto w-full">
                    <table class="w-full text-sm text-left text-gray-500 ">
                        <thead class="text-xs text-gray-700 uppercase bg-gray-50 ">
                            <tr>
                                <th scope="col" class="px-6 py-3">
                                    Node Name
                                </th>
                                <th scope="col" class="px-6 py-3">
                                    User
                                </th>
                                <th scope="col" class="px-6 py-3">
                                    IP
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr class="bg-white" v-for="node in nodes" :key="node.ip">
                                <th scope="row"
                                    class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap flex items-center">
                                    <img class="w-10 h-10 rounded-full mr-2" :src="node.image" alt="user" />
                                    {{ node.name }}
                                </th>
                                <td class="px-6 py-4">
                                    <p>@{{ node.user }}</p>

                                </td>
                                <td class="px-6 py-4">
                                    {{ node.ip }}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>

</template>

<script>

import Verification from '../components/Verification.vue';
import DotLoader from '../components/DotLoader.vue';
import ZerotierApi from "../webapi/api_zerotier.js";

export default {
    name: 'ZerotierPage',
    components: {
        Verification,
        DotLoader
    },
    data() {
        return {
            accesskey: "",
            nodes: [],
            requestedData: false
        }
    },
    computed: {
        Verified() {
            return this.accesskey !== "";
        },
        isLoading() {
            return this.nodes === {} && !this.requestedData;
        },

    },
    methods: {
        OnValidate(accesskey) {
            this.accesskey = accesskey;
            this.getNodes();
        },
        getNodes() {

            ZerotierApi.getNetworkMembers(this.accesskey,
                (data) => {
                    this.nodes = data;
                    this.requestedData = true;
                },
                (err) => {
                    alert("Invalid Key");
                    this.requestedData = true;
                }
            );
        }
    }
}

</script>