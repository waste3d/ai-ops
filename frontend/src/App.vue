<!-- frontend/src/App.vue -->

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { marked } from 'marked';

// --- СОСТОЯНИЕ КОМПОНЕНТА ---
const tickets = ref([]);
const loading = ref(true);
const error = ref(null);

// --- ОСНОВНАЯ ЛОГИКА ---
const fetchTickets = async () => {
  try {
    loading.value = true;
    const response = await axios.get('http://localhost:8000/api/v1/tickets');
    const sortedTickets = response.data.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    tickets.value = sortedTickets;
    error.value = null;
  } catch (err) {
    console.error("Failed to fetch tickets:", err);
    error.value = "Не удалось загрузить инциденты. Запущен ли API Gateway на порту 8000?";
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchTickets();
});

// --- ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ---
const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toLocaleString();
};

const getStatusClass = (status) => {
  if (status === 'analyzed') return 'bg-green-100 text-green-700 border-green-300';
  if (status === 'new') return 'bg-violet-100 text-violet-700 border-violet-300';
  return 'bg-gray-100 text-gray-700 border-gray-300';
};
</script>

<template>
  <div class="min-h-screen font-sans text-gray-800">
    <div class="container mx-auto p-4 md:p-8">
      
      <header class="text-center mb-12 pb-8">
        <h1 class="text-6xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-slate-700 to-violet-700">
          AI-Driven Ops Copilot
        </h1>
        <p class="text-xl text-gray-500 mt-3">Панель мониторинга инцидентов</p>
      </header>

      <div v-if="loading" class="text-center text-gray-500 text-lg">Загрузка...</div>
      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-xl relative text-center shadow-md">{{ error }}</div>
      
      <div v-if="!loading && !error" class="space-y-6">
        <!-- Новая карточка в светлых тонах -->
        <div v-for="ticket in tickets" :key="ticket.id" 
             class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg hover:shadow-2xl hover:shadow-violet-200/50 transition-all duration-300 border border-gray-200/50 overflow-hidden">
          
          <div class="p-6">
            <div class="flex justify-between items-center mb-4">
              <span class="px-3 py-1 text-xs font-bold tracking-wider uppercase rounded-full border" :class="getStatusClass(ticket.status)">
                {{ ticket.status }}
              </span>
              <span class="text-sm text-gray-500 font-medium">Источник: {{ ticket.source }}</span>
            </div>

            <div class="mb-5">
              <p class="text-xl font-semibold text-gray-800">{{ ticket.payload }}</p>
              <div v-if="ticket.analysis_result" class="mt-4 bg-violet-50/50 p-4 rounded-lg border border-violet-200/50">
                <p class="text-sm font-semibold text-violet-800 mb-2">Анализ AI:</p>
                <div class="prose prose-sm max-w-none prose-p:text-gray-700" v-html="marked(ticket.analysis_result)"></div>
              </div>
            </div>

            <div class="flex justify-between items-center text-xs text-gray-400 pt-4 border-t border-gray-200/80">
              <span class="font-mono">ID: {{ ticket.id }}</span>
              <span>{{ formatDate(ticket.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Стили для плагина Typography для светлой темы */
.prose {
  color: #374151; /* text-gray-700 */
}
.prose :where(strong):not(:where([class~="not-prose"] *)) {
  color: #1f2937; /* text-gray-800 */
}
.prose :where(code):not(:where([class~="not-prose"] *)) {
  color: #6d28d9; /* text-violet-700 */
  background-color: #f5f3ff; /* bg-violet-50 */
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
}
.prose :first-child {
  margin-top: 0;
}
.prose :last-child {
  margin-bottom: 0;
}
</style>