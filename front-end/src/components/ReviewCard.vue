<script setup>
import {Star} from 'lucide-vue-next';
import {cn} from "@/lib/utils.js";

const props = defineProps({
  username: String,
  rating: Number,
  review: String,
  class: String,
})
let className = cn("flex flex-col shadow-sm rounded-sm p-5 gap-2", props.class)
</script>

<template>
  <div :class="className">
    <div class="flex items-center">
      <p class="text-lg font-semibold mr-2">{{ username }}</p>
      <Star class="w-5 h-5 text-[#c9a275] fill-[#c9a275]" v-for="index in rating" :key="index"></Star>
    </div>
    <p class="text-balance text-gray-600">{{ review }}</p>
  </div>
</template>