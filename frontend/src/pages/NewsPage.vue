<template>

    <div class="grid h-screen place-items-center">
        
        <Verification v-if="!Verified" :validate="OnValidate" scope="news" />
        <NewsPost v-else :accesskey="accesskey" />

    </div>

</template>

<script>

import Verification from '../components/Verification.vue';
import NewsPost from '../components/NewsPost.vue';

export default {
    name: 'NewsPage',
    components: {
        Verification,
        NewsPost
    },
    data() {
        return {
            accesskey: ""
        }
    },
    computed: {
        Verified() {
            return this.accesskey !== "";
        }
    },
    methods: {
        OnValidate(accesskey) {
            this.accesskey = accesskey;
        }
    }
}

</script>