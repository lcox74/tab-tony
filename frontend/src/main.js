import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router';

import './style.css'

// Pages
import App from './App.vue'
import HomePage from './pages/HomePage.vue';
import NewsPage from './pages/NewsPage.vue';
import ZerotierPage from './pages/ZerotierPage.vue';

// Set Routes
const router = createRouter({
    history: createWebHistory(),
    routes: [
        { 
            path: '/', 
            name: 'Home',
            component: HomePage
        },
        {
            path: '/news',
            name: 'News',
            component: NewsPage 
        },
        { 
            path: '/zerotier', 
            name: 'Zerotier',
            component: ZerotierPage
        },
    ],
});

createApp(App)
.use(router)
.mount('#app')