<template>
    <div
        class="bg-base-200 rounded-[20%] text-4xl flex justify-center items-center p-4"
    >
        <iconify-icon :icon="icon.Name" :style="{ color: icon.Color }">
        </iconify-icon>
    </div>
</template>

<script setup lang="ts">
import type { Icon } from "@/types/service";

defineProps<{ icon: Icon }>();
</script>
