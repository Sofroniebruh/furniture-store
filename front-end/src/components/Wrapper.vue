<script setup>
import {cn} from "@/lib/utils.js";

const props = defineProps({
  class: String,
})
let className = cn("p-5", props.class)
</script>

<template>
  <div :class="className">
    <slot/>
  </div>
</template>