<template>
    <section>
        <h2
            v-if="title"
            class="text-xl tracking-wide mb-2 text-base-content opacity-80 ms-4"
        >
            {{ title }}
        </h2>
        <div class="flex flex-wrap gap-4">
            <ServiceCard
                v-for="service in services"
                :key="service.Link"
                :service="service"
                :tags="tags"
            />
        </div>
    </section>
</template>

<script setup lang="ts">
import ServiceCard from "@/components/Service.vue";
import type { Service, Tag } from "@/types/service";

defineProps<{
    title?: string;
    services: Service[];
    tags: Tag[];
}>();
</script>
