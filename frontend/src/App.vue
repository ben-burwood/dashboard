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

// Resolve a favicon config value to an href.
// Iconify name becomes an Iconify SVG API URL; a URL / path / data: URI is used verbatim.
const faviconHref = (value?: string): string | null => {
    const v = (value ?? "").trim();
    if (!v) return null;
    if (
        v.includes("://") ||
        v.startsWith("/") ||
        v.startsWith("./") ||
        v.startsWith("../") ||
        v.startsWith("data:")
    ) {
        return v;
    }
    const sep = v.indexOf(":");
    if (sep !== -1) {
        return `https://api.iconify.design/${v.slice(0, sep)}/${v.slice(sep + 1)}.svg`;
    }
    return v;
};

const applyFavicon = (value?: string) => {
    const href = faviconHref(value);
    if (!href) return;
    let link = document.querySelector<HTMLLinkElement>('link[rel="icon"]');
    if (!link) {
        link = document.createElement("link");
        link.rel = "icon";
        document.head.appendChild(link);
    }
    link.href = href;
};

onMounted(async () => {
    try {
        const res = await fetch(`${SERVER_URL}/config`);
        const data: Config = await res.json();
        config.value = data;
        applyFavicon(data.Favicon);
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
