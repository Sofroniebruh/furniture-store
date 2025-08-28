<script setup lang="ts">

import {Button} from "@/components/ui/button/index.js";
import {Info, Trash} from 'lucide-vue-next';
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.js";
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory.js";
import {PropType} from "vue";
import {Product} from "@/lib/types";

const {toggleWishlist, addToWishlistStore} = wishlistStore
const {removeFromWishlistComposable} = useWishlistOrHistory()
const {setAllSheetsClosed} = useScreenSheetStore()

defineProps({
  product: {
    type: Object as PropType<Product>
  },
  isHistoryCard: {
    type: Boolean,
    default: false
  }
})

const handleRemoveFromWishlist = async (product: Product) => {
  await toggleWishlist(product)

  try {
    const res = await removeFromWishlistComposable(product.id)

    if (res !== 200) {
      addToWishlistStore(product)
    }
  } catch (e) {
    console.error(e)
    addToWishlistStore(product)
  }
}
</script>

<template>
  <div class="flex justify-between items-center">
    <div class="flex items-center gap-2">
      <div class="overflow-hidden w-[100px] h-auto">
        <img class="object-cover w-full h-full" :src="product.pictureUrls[0]" alt="product image">
      </div>
      <div>
        <p class="text-base text-gray-700">{{ product.name }}</p>
        <p class="text-sm font-semibold">{{ product.price }} &#8364;</p>
      </div>
    </div>
    <Button v-if="!isHistoryCard" @click="() => handleRemoveFromWishlist(product)" size="sm" class="cursor-pointer"
            variant="outline">
      <Trash/>
    </Button>
    <router-link v-else @click="setAllSheetsClosed" :to="`/product/${product.id}`">
      <Button size="sm" class="cursor-pointer" variant="outline">
        <Info/>
      </Button>
    </router-link>
  </div>
</template>