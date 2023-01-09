<template>

    <div class="flex justify-center my-auto">
        <div class="block p-6 rounded-lg shadow-lg bg-white max-w-sm">
            <h5 class="text-gray-900 text-xl leading-tight font-medium mb-2">Welcome to the TAB news form</h5>
            <p class="text-gray-700 text-base mb-4">
                Please enter your access key below to post tech news, this is to keep only memebers of the TAB server
                from posting.
                <br><br>
                If you are apart of the server but don't have an access key, type <code
                    class="bg-gray-300 rounded px-1">/request news</code> in the
                server's general chat.
            </p>

            <div class="flex mb-3">
                <div class="flex-grow">
                    <input
                        class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
                        type="text" placeholder="Access Key..." :value="accesskey" @input="e => accesskey = e.target.value" >
                </div>
                <div class="flex-none">
                    <button
                        class="mx-2 fill-current bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" @click="checkValid">
                        Verify
                    </button>

                </div>
            </div>



        </div>
    </div>

</template>
<script>

import NewsApi from "../webapi/api_news.js";

export default {
    name: 'Verification',
    props: {
        validate: {
            type: Function,
            required: true
        }
    },
    data() {
        return {
            accesskey: ""
        }
    },
    methods: {
        checkValid() {
            // Make the API request
            NewsApi.validateAccess(this.accesskey,
                (data) => {
                    this.validate(this.accesskey);
                    this.accesskey = "";
                },
                (err) => {
                    alert("Invalid Key");
                }
            );

        }
    }
}

</script>