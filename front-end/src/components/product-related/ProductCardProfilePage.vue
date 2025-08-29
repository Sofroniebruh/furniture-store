<script setup lang="ts">
import {PropType} from "vue";
import {Product} from "@/lib/types";
import WishlistDelete from "@/components/product-related/WishlistDelete.vue";
import PurchaseHistory from "@/components/product-related/PurchaseHistory.vue";

const props = defineProps({
  product: {
    type: Object as PropType<Product>
  },
  isHistoryCard: {
    default: false
  }
})
const product = props.product;
</script>

<template>
  <div class="flex flex-col">
    <div class="rounded-sm w-[150px] h-[150px] overflow-hidden">
      <img class="w-full h-full object-cover" :src="product.pictureUrls[0]" alt="product image">
    </div>
    <div class="flex items-center justify-between">
      <p>{{ product.name }}</p>
      <WishlistDelete v-if="!isHistoryCard" :product="product"></WishlistDelete>
      <PurchaseHistory v-else :product="product"/>
    </div>
  </div>
</template>