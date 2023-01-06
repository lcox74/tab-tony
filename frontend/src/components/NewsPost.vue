<template>

    <div v-if="accesskey !== ''" class="flex justify-center my-auto">
        <div class="block p-6 rounded-lg shadow-lg bg-white max-w-md">
            <h5 class="text-gray-900 text-xl leading-tight font-medium mb-2">Create a News Post</h5>
            <p class="text-gray-700 text-base mb-4">Please keep it tech related, this can be from software, companies,
                automotive, hardware etc.</p>


            <div class="flex">
                <span :class="style.prependInputStyle">#</span>
                <input type="text" :class="style.inputStyle" class="rounded-none rounded-r-lg" placeholder="Post Title"
                    :value="form.title" @input="e => form.title = e.target.value" />
            </div>

            <textarea rows="5" :class="style.inputStyle" class="mt-2" placeholder="Post Content..."
                :value="form.content" @input="e => form.content = e.target.value" />

            <input type="text" :class="style.inputStyle" class="mt-2" placeholder="News Source URL" :value="form.url"
                @input="e => form.url = e.target.value" />
            <input type="text" :class="style.inputStyle" class="mt-2" placeholder="Image URL [Optional]"
                :value="form.image_url" @input="e => form.image_url = e.target.value" />

            <button @click="createPost"
                class="mt-2 float-right bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded">
                Post
            </button>

        </div>
    </div>


</template>
<script>

export default {
    name: 'NewsPost',
    props: {
        accesskey: {
            type: String,
            required: true
        }
    },
    data() {
        return {
            style: {
                inputStyle: "bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5",
                prependInputStyle: "inline-flex items-center px-3 text-sm text-gray-600 font-semibold bg-gray-200 border border-r-0 border-gray-300 rounded-l-md",
            },
            form: {
                title: "",
                content: "",
                tags: [],
                image_url: "",
                url: "",
                access: "",
                clean() {
                    this.title = "";
                    this.content = "";
                    this.tags = [];
                    this.image_url = "";
                    this.url = "";
                    this.access = "";
                }
            }
        }
    },
    methods: {
        createPost() {
            this.form.access = this.accesskey;
            if (this.form.title === "" || this.form.content === "" || this.form.url === "") {
                alert("Please fill out all fields");
                return;
            }

            if (this.form.title[0] != "#") {
                this.form.title = "#" + this.form.title;
            }

            fetch("http://localhost:3000/news", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(this.form)
            }).then(res => res.json()).then(data => {
                console.log(data);
                this.form.clean();
            }).catch(err => {
                console.log(err);
            })

        }
    }
}

</script>