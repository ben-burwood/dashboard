<template>
    <span
        class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
        :class="getBadgeClass(tag.Color)"
    >
        {{ tag.Name }}
    </span>
</template>

<script setup lang="ts">
import type { Tag } from "@/types/service";

defineProps<{ tag: Tag }>();

const colorVariants = {
    primary: "bg-indigo-100 text-indigo-700 dark:bg-indigo-500/20 dark:text-indigo-300",
    secondary: "bg-purple-100 text-purple-700 dark:bg-purple-500/20 dark:text-purple-300",
    accent: "bg-teal-100 text-teal-700 dark:bg-teal-500/20 dark:text-teal-300",
    info: "bg-sky-100 text-sky-700 dark:bg-sky-500/20 dark:text-sky-300",
    success: "bg-green-100 text-green-700 dark:bg-green-500/20 dark:text-green-300",
    warning: "bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300",
    error: "bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-300",
    neutral: "bg-gray-200 text-gray-700 dark:bg-gray-600/40 dark:text-gray-300",
};

const getBadgeClass = (color: string) =>
    colorVariants[color as keyof typeof colorVariants] ?? colorVariants.neutral;
</script>
