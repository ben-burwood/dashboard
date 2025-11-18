<template>
    <a :href="service.Link" class="transition-transform hover:scale-[1.025]">
        <div
            class="flex items-center gap-4 rounded-2xl bg-base-100 opacity-90 shadow-sm hover:shadow-lg p-4 min-w-[260px] max-w-md"
        >
            <ServiceIcon :icon="service.Icon" />
            <div class="flex-1">
                <h2 class="text-lg font-semibold mb-2">{{ service.Title }}</h2>
                <span class="flex flex-wrap gap-2">
                    <ServiceTag
                        :tag="tag"
                        v-for="tag in serviceTags"
                        :key="tag.Name"
                    />
                </span>
            </div>
        </div>
    </a>
</template>

<script setup lang="ts">
import ServiceIcon from "@/components/ServiceIcon.vue";
import ServiceTag from "@/components/ServiceTag.vue";
import type { Service, Tag } from "@/types/service";

const props = defineProps<{
    service: Service;
    tags: Array<Tag>;
}>();

const serviceTags = !props.service.Tags
    ? []
    : props.service.Tags.map((name) =>
          props.tags.find((tag) => tag.Name === name),
      ).filter((tag): tag is Tag => !!tag);
</script>
