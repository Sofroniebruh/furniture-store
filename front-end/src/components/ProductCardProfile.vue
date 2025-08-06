<script setup>

import {Button} from "@/components/ui/button/index.js";
import {Info, Trash} from 'lucide-vue-next';

import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory.js";

const {removeFromWishlist} = useWishlistOrHistory()

defineProps({
  imageSrc: String,
  name: String,
  price: Number,
  id: String,
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
        <img class="object-cover w-full h-full" :src="imageSrc" alt="product image">
      </div>
      <div>
        <p class="text-base text-gray-700">{{ name }}</p>
        <p class="text-sm font-semibold">{{ price }} &#8364;</p>
      </div>
    </div>
    <Button v-if="!isHistoryCard" @click="removeFromWishlist(id)" size="sm" class="cursor-pointer" variant="outline">
      <Trash/>
    </Button>
    <Button v-if="isHistoryCard" size="sm" class="cursor-pointer" variant="outline">
      <Info/>
    </Button>
  </div>
</template>