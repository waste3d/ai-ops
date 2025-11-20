import { createApp } from 'vue'
import './index.css' // <-- УБЕДИТЕСЬ, ЧТО ЗДЕСЬ ИМЕННО 'index.css'
import App from './App.vue'
import Particles from "vue3-particles";
import { loadSlim } from "tsparticles-slim";

const app = createApp(App);

app.use(Particles, {
    init: async engine => {
        await loadSlim(engine); 
    },
});

app.mount('#app');