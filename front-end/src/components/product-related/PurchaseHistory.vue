<script setup lang="ts">

import {Info} from "lucide-vue-next";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore";
import {PropType} from "vue";
import {Product} from "@/lib/types";

const {setAllSheetsClosed} = useScreenSheetStore()

const props = defineProps({
  product: {
    type: Object as PropType<Product>
  },
})
const product = props.product;
</script>

<template>
  <router-link class="h-[24px]" @click="setAllSheetsClosed" :to="`/product/${product.id}`">
    <Button size="sm" class="cursor-pointer" variant="outline">
      <Info/>
    </Button>
  </router-link>
</template>