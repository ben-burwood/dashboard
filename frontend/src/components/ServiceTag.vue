<template>
    <span class="badge badge-soft" :class="getBadgeClass(tag.Color)">
        {{ tag.Name }}
    </span>
</template>

<script setup lang="ts">
import type { Tag } from '@/types/service';

defineProps<{ tag: Tag }>();

const colorVariants = {
    primary: "badge-primary",
    secondary: "badge-secondary",
    accent: "badge-accent",
    info: "badge-info",
    success: "badge-success",
    warning: "badge-warning",
    error: "badge-error",
    neutral: "badge-neutral",
}

const getBadgeClass = (color: string) => colorVariants[color as keyof typeof colorVariants];
</script>
