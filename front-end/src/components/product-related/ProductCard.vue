<script setup lang="ts">
import Wrapper from "@/components/Wrapper.vue";
import {Button} from "@/components/ui/button";
import {Plus} from 'lucide-vue-next';
import {cn} from "@/lib/utils.js";
import {Product} from "@/lib/types";
import {PropType} from "vue";

const props = defineProps({
  product: {
    type: Object as PropType<Product>,
    required: true
  },
  class: {
    type: String,
    required: false,
  }
})

const className = props.class
const product = props.product
</script>

<template>
  <div :class="cn('w-full h-full', className)">
    <Wrapper class="p-3 sm:p-4 h-full">
      <div class="flex flex-col gap-3 h-full">
        <div class="relative w-full aspect-square overflow-hidden rounded-lg bg-gray-100">
          <img
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
              :src="product.pictureUrls[0]"
              :alt="`Image of ${product.name}`"
              loading="lazy"
          />
          <div v-if="product.event && product.event !== 'none'"
               class="absolute top-2 left-2 px-2 py-1 text-xs font-medium text-white bg-red-500 rounded-full capitalize">
            {{ product.event }}
          </div>
        </div>

        <div class="flex-1 flex flex-col gap-2">
          <div class="flex-1">
            <h3 class="text-sm sm:text-base font-medium text-gray-900 line-clamp-2 leading-tight">
              {{ product.name }}
            </h3>
            <div class="flex justify-between items-center mt-2">
              <p class="text-lg font-bold text-gray-900">{{ product.price }}€</p>
              <div v-if="product.colors?.length"
                   class="px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded-full">
                {{ product.colors.length }} {{ product.colors.length === 1 ? 'color' : 'colors' }}
              </div>
            </div>
          </div>
          <router-link :to="`/product/${product.id}`">
            <Button
                class="w-full cursor-pointer bg-[#c9a275] hover:bg-[#dbb384] text-white transition-colors duration-200 text-sm"
            >
              Add to cart
              <Plus class="w-4 h-4 ml-1"/>
            </Button>
          </router-link>
        </div>
      </div>
    </Wrapper>
  </div>
</template>