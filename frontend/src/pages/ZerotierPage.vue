<template>

    <div class="grid h-screen place-items-center">
        
        <Verification v-if="!Verified" :validate="OnValidate" scope="zerotier" />
        <div v-else>
            <div v-if="isLoading">
                <DotLoader />
            </div>
            <div v-else>
                <div v-for="user in Object.keys(nodes)" :key="user">
                    <ZerotierUser :user="user" :nodes="nodes[user]" />
                </div>
            </div>
        </div>
    </div>

</template>

<script>

import Verification from '../components/Verification.vue';
import ZerotierUser from '../components/ZerotierUser.vue';
import DotLoader from '../components/DotLoader.vue';

export default {
    name: 'ZerotierPage',
    components: {
        Verification,
        ZerotierUser,
        DotLoader
    },
    data() {
        return {
            accesskey: "",
            nodes: {},
            requestedData: false
        }
    },
    computed: {
        Verified() {
            return this.accesskey !== "";
        },
        isLoading() {
            return this.nodes === {} && !this.requestedData;
        }
    },
    methods: {
        OnValidate(accesskey) {
            this.accesskey = accesskey;
            this.getNodes();
        },
        getNodes() {
            ZerotierApi.getNodes(this.accesskey,
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