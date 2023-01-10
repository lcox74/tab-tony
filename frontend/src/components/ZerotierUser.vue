<template>
    <div class="max-w-sm w-full lg:max-w-full lg:flex">
        <div class="h-48 lg:h-auto lg:w-48 flex-none bg-cover rounded-t lg:rounded-t-none lg:rounded-l text-center overflow-hidden"
            :style="{ backgroundImage: getAvatarImgUrl() }" :title="user">
        </div>
        <div
            class="border-r border-b border-l border-gray-400 lg:border-l-0 lg:border-t lg:border-gray-400 bg-white rounded-b lg:rounded-b-none lg:rounded-r p-4 flex flex-col justify-between leading-normal">
            <div class="mb-8">
                <div class="text-gray-900 font-bold text-xl mb-2">{{ user }}'s Nodes</div>
                <div class="flex">
                    <div v-for="node in nodes">
                        <div class="text-gray-700 text-base">{{ node.name }}</div>
                        <div class="text-gray-700 text-base">{{ node.ip }}</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'ZerotierUser',
    props: {
        user: {
            type: String,
            required: true
        },
        nodes: {
            type: Array,
            required: true
        }
    },
    methods: {
        getAvatarImgUrl() {
            return this.nodes[0].image;
        }
    }
}
</script>