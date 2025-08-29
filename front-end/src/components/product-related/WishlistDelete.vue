<script setup lang="ts">
import {useWishlistStore} from "@/stores/useWishlist.js";
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory.js";
import {PropType} from "vue";
import {Product} from "@/lib/types.js";
import {Trash} from "lucide-vue-next";

const wishlistStore = useWishlistStore()
const {toggleWishlist, addToWishlistStore} = wishlistStore
const {removeFromWishlistComposable} = useWishlistOrHistory()

defineProps({
  product: {
    type: Object as PropType<Product>
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
  <Button @click="() => handleRemoveFromWishlist(product)" size="sm" class="cursor-pointer"
          variant="outline">
    <Trash/>
  </Button>
</template>