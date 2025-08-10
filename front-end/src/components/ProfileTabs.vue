<script setup lang="ts">
import {Tabs, TabsContent, TabsList, TabsTrigger,} from '@/components/ui/tabs'
import {useWishlistOrHistory} from "@/composables/useWishlistOrHistory";
import {onMounted} from "vue";
import ProductCardProfile from "@/components/ProductCardProfile.vue";
import {useWishlistStore} from "@/stores/useWishlist";

const {
  productsHistoryData,
  fetchWishlistData,
  wishlistLoading,
  wishlistError,
  historyLoading,
  historyError
} = useWishlistOrHistory()
const {itemsInWishlist} = useWishlistStore()

onMounted(async () => {
  await fetchWishlistData("wishlist")
  await fetchWishlistData("history")
})
</script>

<template>
  <div class="h-screen flex flex-col">
    <Tabs default-value="wishlist" class="w-[350px] flex-1 flex flex-col">
      <TabsList class="grid w-full grid-cols-2 flex-shrink-0">
        <TabsTrigger value="wishlist">
          Wishlist
        </TabsTrigger>
        <TabsTrigger value="history">
          Purchase history
        </TabsTrigger>
      </TabsList>
      <TabsContent value="wishlist" class="flex-1 flex flex-col">
        <div class="h-[calc(100%-185px)] overflow-y-auto flex flex-col">
          <div v-if="wishlistError || wishlistLoading || itemsInWishlist.length == 0"
               class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600" v-if="wishlistLoading">Loading...</p>
            <p class="text-base font-semibold text-gray-600" v-if="wishlistError">Error on our side...</p>
            <p class="text-base font-semibold text-gray-600"
               v-if="!wishlistLoading && itemsInWishlist.length == 0">No products
              added</p>
          </div>
          <ProductCardProfile v-for="(product, index) in itemsInWishlist" :image-src="product.pictureUrls[0]"
                              :name="product.name" :price="product.price" :id="product.id" :product="product"
                              :key="index"/>
        </div>
      </TabsContent>
      <TabsContent value="history" class="flex-1 flex flex-col">
        <div class="h-[calc(100%-185px)] overflow-y-auto flex flex-col">
          <div v-if="historyError || historyLoading || productsHistoryData.length == 0"
               class="w-full h-full flex flex-col justify-center items-center">
            <p class="text-base font-semibold text-gray-600" v-if="historyLoading">Loading...</p>
            <p class="text-base font-semibold text-gray-600" v-if="historyError">Error on our side...</p>
            <p class="text-base font-semibold text-gray-600" v-if="productsHistoryData.length == 0">No purchases
              made</p>
          </div>
          <ProductCardProfile :is-history-card="true" v-for="(product, index) in productsHistoryData"
                              :image-src="product.pictureUrls[0]"
                              :product="product"
                              :name="product.name" :price="product.price" :id="product.id" :key="index"/>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>