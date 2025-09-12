<script setup>

import ProductCard from "@/components/product-related/ProductCard.vue";
import Wrapper from "@/components/Wrapper.vue";
import Sorting from "@/components/Sorting.vue";
import {SlidersHorizontal} from 'lucide-vue-next';
import {Button} from "@/components/ui/button/index.js";
import SheetGeneral from "@/components/SheetGeneral.vue";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore.ts";
import {useProducts} from "@/composables/useProducts.js";
import {onMounted} from "vue";
import Pagination from "@/components/Pagination.vue";

const smallScreenSheet = useScreenSheetStore()
const {
  products,
  error,
  fetchProducts,
  currentPage,
  handleNext,
  setPage,
  handlePrevious,
  totalPages
} = useProducts()

onMounted(async () => {
  await fetchProducts()
})

</script>

<template>
  <div class="flex flex-col md:flex-row h-full">
    <Wrapper class="md:hidden block">
      <SheetGeneral :is-open="smallScreenSheet.isSheetOpen('SmallScreenFiltering')" title="FUMI" side="left">
        <template #trigger>
          <Button @click="smallScreenSheet.setOpenSheet(true, {name: 'SmallScreenFiltering'})" class="cursor-pointer"
                  variant="outline">Filters
            <SlidersHorizontal/>
          </Button>
        </template>
        <template #content>
          <Wrapper class="pt-0">
            <Sorting/>
          </Wrapper>
        </template>
      </SheetGeneral>
    </Wrapper>
    <Wrapper class="hidden w-1/4 md:block pr-0">
      <Sorting/>
    </Wrapper>
    <Wrapper
        class="p-0 md:pl-5 md:pt-5 pb-0 w-full md:w-3/4 flex flex-col items-center min-h-[calc(100dvh-200px-68px-76px)]">
      <div v-if="error === null && products === null"
           class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6 w-full">
        <div v-for="n in 12" :key="n" 
             class="bg-gray-200 animate-pulse rounded-lg aspect-[3/4] w-full"></div>
      </div>
      <div v-else-if="products === undefined || products === null || products.length === 0"
           class="flex justify-center items-center h-[calc(100vh-300px)]">
        <div class="text-center">
          <h1 class="font-semibold text-base sm:text-lg mb-2">No products found</h1>
          <p class="text-gray-600 text-sm">Try adjusting your filters or search terms</p>
        </div>
      </div>
      <div v-else-if="error" class="text-center py-8">
        <p class="text-red-600 mb-4">Failed to load products. Please try again.</p>
        <Button @click="fetchProducts" class="mt-4">Retry</Button>
      </div>
      <div v-else
           class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6 w-full"
           role="grid" 
           aria-label="Product listings">
        <router-link 
          v-for="(product, index) in products" 
          :key="product.id"
          :to="`/product/${product.id}`"
          class="group focus:outline-none focus-visible:ring-2 focus-visible:ring-[#c9a275] rounded-lg"
          role="gridcell"
          :aria-label="`View ${product.name} details`">
          <ProductCard 
            :product="product" 
            class="w-full" />
        </router-link>
      </div>
      <Pagination :set-page="setPage"
                  :current-page="currentPage"
                  :handle-next="handleNext"
                  :handle-previous="handlePrevious"
                  :total-pages="totalPages"
                  v-if="products !== null"/>
    </Wrapper>
  </div>
</template>