<template>
    <div class="min-h-screen bg-gray-200 text-gray-900 dark:bg-gray-900 dark:text-gray-100">
        <h2
            v-if="config?.Title"
            class="text-4xl font-semibold text-center opacity-60 tracking-wide py-4"
        >
            {{ config.Title }}
        </h2>
        <div class="mx-auto flex max-w-7xl flex-col items-start gap-10 px-4 py-8">
            <ServiceGroup
                v-for="(groupServices, group) in grouped"
                :key="group"
                :title="group"
                :services="groupServices"
                :tags="config?.Tags ?? []"
            />
        </div>
        <ThemeSwitcher class="fixed bottom-4 right-4" />
    </div>
</template>

<script setup lang="ts">
import ServiceGroup from "@/components/ServiceGroup.vue";
import { SERVER_URL } from "@/main";
import type { Config } from "@/types/config";
import type { Service } from "@/types/service";
import { computed, onMounted, ref } from "vue";
import ThemeSwitcher from "./components/ThemeSwitcher.vue";

const config = ref<Config | null>(null);

onMounted(async () => {
    try {
        const res = await fetch(`${SERVER_URL}/config`);
        const data: Config = await res.json();
        config.value = data;
    } catch (e) {
        console.error("Failed to fetch config:", e);
        config.value = null;
    }
});

const DEFAULT_GROUP = "default";

const grouped = computed(() => {
    const result: Record<string, Service[]> = {};
    for (const service of config.value?.Services || []) {
        const group = service.Group || DEFAULT_GROUP;
        if (!result[group]) result[group] = [];
        result[group].push(service);
    }
    return result;
});
</script>
