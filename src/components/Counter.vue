<template>
  <div class="inline-flex items-center gap-4 px-4 py-2 bg-primary-100 dark:bg-primary-900 rounded-lg">
    <button
      @click="count--"
      class="w-8 h-8 flex items-center justify-center rounded-full bg-primary-500 text-white hover:bg-primary-600 transition-colors"
    >
      -
    </button>
    <span class="text-2xl font-bold min-w-[2rem] text-center">{{ count }}</span>
    <button
      @click="count++"
      class="w-8 h-8 flex items-center justify-center rounded-full bg-primary-500 text-white hover:bg-primary-600 transition-colors"
    >
      +
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const count = ref(0);
</script>