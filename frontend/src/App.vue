<template>
  <div class="min-h-screen bg-gray-100">
  <h2 v-if="config?.Title" class="text-4xl font-semibold text-center text-gray-600 tracking-wide py-4">{{ config.Title }}</h2>
    <div class="hero">
      <div class="hero-content flex flex-col items-start gap-10">
        <ServiceGroup
        v-for="(groupServices, group) in grouped"
        :key="group"
        :title="group"
        :services="groupServices"
        :tags="config?.Tags ?? []" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ServiceGroup from '@/components/ServiceGroup.vue';
import { SERVER_URL } from '@/main';
import type { Config } from '@/types/config';
import type { Service } from '@/types/service';
import { computed, onMounted, ref } from 'vue';

const config = ref<Config | null>(null);

onMounted(async () => {
  try {
    const res = await fetch(`${SERVER_URL}/config`);
    const data: Config = await res.json();
    config.value = data;
  } catch (e) {
    console.error('Failed to fetch config:', e);
    config.value = null;
  }
});

const DEFAULT_GROUP = 'default'

const grouped = computed(() => {
  const result: Record<string, Service[]> = {}
  for (const service of config.value?.Services || []) {
    const group = service.Group || DEFAULT_GROUP
    if (!result[group]) result[group] = []
    result[group].push(service)
  }
  return result
})
</script>
