<script setup lang="ts">

import {PropType} from "vue";
import {Product} from "@/lib/types";
import WishlistDelete from "@/components/product-related/WishlistDelete.vue";
import PurchaseHistory from "@/components/product-related/PurchaseHistory.vue";


defineProps({
  product: {
    type: Object as PropType<Product>
  },
  isHistoryCard: {
    type: Boolean,
    default: false
  }
})
</script>

<template>
  <div class="flex justify-between items-center">
    <div class="flex items-center gap-2">
      <div class="overflow-hidden w-[100px] h-auto">
        <img class="object-cover w-full h-full" :src="product!.pictureUrls[0]" alt="product image">
      </div>
      <div>
        <p class="text-base text-gray-700">{{ product!.name }}</p>
        <p class="text-sm font-semibold">{{ product!.price }} &#8364;</p>
      </div>
    </div>
    <WishlistDelete :product="product" v-if="!isHistoryCard"/>
    <PurchaseHistory :product="product" v-else/>
  </div>
</template>