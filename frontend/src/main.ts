import { createApp } from 'vue';
import App from './App.vue';
import './style.css';

const savedTheme = localStorage.getItem('theme');
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
if (savedTheme === 'dark' || (!savedTheme && prefersDark)) {
    document.documentElement.classList.add('dark');
}

export const SERVER_URL = `${window.location.protocol}//${window.location.hostname}:${window.location.port}/api`;

createApp(App).mount('#app')
