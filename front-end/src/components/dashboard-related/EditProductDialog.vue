<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Product, Color } from '@/lib/types'

interface Props {
  product: Product | null
  isOpen: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
  'product-updated': [product: Product]
}>()

const form = ref({
  name: '',
  description: '',
  price: 0,
  stock: 0,
  model: '',
  event: 'none',
  colors: [] as Color[]
})

const newColor = ref('')
const isSubmitting = ref(false)

watch(() => props.product, (product) => {
  if (product) {
    form.value = {
      name: product.name,
      description: product.description,
      price: product.price,
      stock: product.stock,
      model: product.model,
      event: product.event,
      colors: [...product.colors]
    }
  }
}, { immediate: true })

const addColor = () => {
  if (newColor.value.trim()) {
    form.value.colors.push({
      name: newColor.value.trim()
    })
    newColor.value = ''
  }
}

const removeColor = (index: number) => {
  form.value.colors.splice(index, 1)
}

const close = () => {
  emit('close')
}

const updateProduct = async () => {
  if (!props.product) return
  
  isSubmitting.value = true
  
  try {
    const formData = new FormData()
    formData.append('name', form.value.name)
    formData.append('description', form.value.description)
    formData.append('price', form.value.price.toString())
    formData.append('stock', form.value.stock.toString())
    formData.append('model', form.value.model.toLowerCase())
    formData.append('event', form.value.event)
    
    const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products?id=${props.product.id}`, {
      method: 'PUT',
      credentials: 'include',
      body: formData
    })
    
    if (!response.ok) {
      throw new Error(`Failed to update product: ${response.status}`)
    }
    
    const data = await response.json()
    emit('product-updated', data.updated)
    close()
  } catch (error) {
    console.error('Failed to update product:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="w-full max-w-2xl max-h-[90vh] overflow-y-auto">
    <form @submit.prevent="updateProduct" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Product Name</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
          >
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea
            v-model="form.description"
            required
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
          ></textarea>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Price ($)</label>
            <input
              v-model.number="form.price"
              type="number"
              step="0.01"
              min="0"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Stock Quantity</label>
            <input
              v-model.number="form.stock"
              type="number"
              min="0"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
            >
          </div>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Model</label>
            <select
                v-model="form.model"
                style="-webkit-appearance: none;"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
            >
              <option value="" disabled>Select a model</option>
              <option value="bed">Bed</option>
              <option value="chair">Chair</option>
              <option value="table">Table</option>
              <option value="sofa">Sofa</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Event</label>
            <select
              v-model="form.event"
              style="-webkit-appearance: none;"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
            >
              <option value="none">None</option>
              <option value="sale">Sale</option>
              <option value="new">New</option>
              <option value="featured">Featured</option>
              <option value="bestseller">Bestseller</option>
            </select>
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Colors</label>
          <div class="space-y-2">
            <div class="flex flex-wrap gap-2 mb-2">
              <span
                v-for="(color, index) in form.colors"
                :key="index"
                class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
              >
                {{ color.name }}
                <button
                  type="button"
                  @click="removeColor(index)"
                  class="ml-1 text-gray-400 hover:text-gray-600"
                >
                  ×
                </button>
              </span>
            </div>
            <div class="flex gap-2">
              <input
                v-model="newColor"
                type="text"
                placeholder="Add color..."
                class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#c9a275] focus:border-transparent"
                @keyup.enter.prevent="addColor"
              >
              <button
                type="button"
                @click="addColor"
                class="px-4 py-2 bg-gray-100 cursor-pointer text-gray-700 rounded-lg hover:bg-gray-200"
              >
                Add
              </button>
            </div>
          </div>
        </div>
        <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 pt-4 border-t">
          <button
            type="button"
            @click="close"
            class="flex-1 px-4 py-2 border cursor-pointer border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex-1 px-4 py-2 bg-[#c9a275] text-white cursor-pointer rounded-lg hover:bg-[#b8956a] font-medium disabled:opacity-50"
          >
            {{ isSubmitting ? 'Updating...' : 'Update Product' }}
          </button>
        </div>
      </form>
  </div>
</template>