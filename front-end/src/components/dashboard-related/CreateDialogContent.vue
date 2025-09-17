<script setup lang="ts">
import {computed, onMounted, reactive, ref} from 'vue'
import type {Color, Product} from '@/lib/types'
import {useScreenSheetStore} from "@/stores/useScreenSheetStore"
import {useCreateProduct} from "@/composables/useCreateProduct";
import {cn} from "@/lib/utils.js";
import {toast} from "vue-sonner";

let allColors = ref<Color[]>([])
const {addProduct, fetchColors} = useCreateProduct()

onMounted(async () => {
  await fetchColors().then((res) => {
    allColors.value = res.colors!
  })
})

const formData = reactive({
  name: '',
  description: '',
  price: 0,
  amount: 0,
  model: '',
  event: 'none',
  images: [] as File[],
  colors: [{name: ''}] as Color[],
})

const imagePreviewUrls = ref<string[]>([])

const errors = reactive({
  name: '',
  description: '',
  price: '',
  amount: '',
  model: '',
  images: '',
  colors: '',
  general: ''
})

const isSubmitting = ref(false)

const {setOpenDialog} = useScreenSheetStore()

const validateForm = (): boolean => {
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  if (!formData.name.trim()) {
    errors.name = 'Product name is required'
    isValid = false
  } else if (formData.name.trim().length < 2) {
    errors.name = 'Product name must be at least 2 characters'
    isValid = false
  } else if (formData.name.trim().length > 100) {
    errors.name = 'Product name must be less than 100 characters'
    isValid = false
  }

  if (!formData.description.trim()) {
    errors.description = 'Product description is required'
    isValid = false
  } else if (formData.description.trim().length < 10) {
    errors.description = 'Description must be at least 10 characters'
    isValid = false
  } else if (formData.description.trim().length > 1000) {
    errors.description = 'Description must be less than 1000 characters'
    isValid = false
  }

  if (!formData.model.trim()) {
    errors.model = 'Product model is required'
    isValid = false
  } else if (formData.model.trim().length > 50) {
    errors.model = 'Model must be less than 50 characters'
    isValid = false
  }

  if (formData.price <= 0) {
    errors.price = 'Price must be greater than 0'
    isValid = false
  } else if (formData.price > 999999.99) {
    errors.price = 'Price must be less than $999,999.99'
    isValid = false
  }

  if (formData.amount < 0) {
    errors.amount = 'Stock amount cannot be negative'
    isValid = false
  } else if (formData.amount > 999999) {
    errors.amount = 'Stock amount must be less than 999,999'
    isValid = false
  }

  if (formData.images.length === 0) {
    errors.images = 'At least one image is required'
    isValid = false
  } else {
    const maxSize = 5 * 1024 * 1024
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']

    const hasInvalidFile = formData.images.some(file => {
      return !allowedTypes.includes(file.type) || file.size > maxSize
    })

    if (hasInvalidFile) {
      errors.images = 'Images must be JPEG, PNG, WebP, or GIF format and less than 5MB'
      isValid = false
    }
  }

  const validColors = formData.colors.filter(color => color.name.trim())
  if (validColors.length === 0) {
    errors.colors = 'At least one color is required'
    isValid = false
  }

  return isValid
}

const addColor = (color?: Color) => {
  formData.colors.push(color ? {id: color.id, name: color.name} : {id: undefined, name: ""})
}

const isColorChosen = (name: string) => {
  return formData.colors.some(color => color.name === name)
}

const removeColor = (index: number) => {
  if (formData.colors.length > 1) {
    formData.colors.splice(index, 1)
  }
}

const addPicture = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.multiple = true
  input.onchange = (e) => {
    const files = (e.target as HTMLInputElement).files
    if (files) {
      handleImageUpload(Array.from(files))
    }
  }
  input.click()
}

const removePictureUrl = (index: number) => {
  formData.images.splice(index, 1)
  if (imagePreviewUrls.value[index]) {
    URL.revokeObjectURL(imagePreviewUrls.value[index])
  }
  imagePreviewUrls.value.splice(index, 1)
}

const handleImageUpload = (files: File[]) => {
  const maxSize = 5 * 1024 * 1024
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']

  const validFiles = files.filter(file => {
    return allowedTypes.includes(file.type) && file.size <= maxSize
  })

  if (validFiles.length !== files.length) {
    errors.images = 'Some files were skipped. Only JPEG, PNG, WebP, and GIF files under 5MB are allowed.'
  } else {
    errors.images = ''
  }

  formData.images.push(...validFiles)

  validFiles.forEach(file => {
    const previewUrl = URL.createObjectURL(file)
    imagePreviewUrls.value.push(previewUrl)
  })
}

const resetForm = () => {
  imagePreviewUrls.value.forEach(url => URL.revokeObjectURL(url))

  formData.name = ''
  formData.description = ''
  formData.price = 0
  formData.amount = 0
  formData.model = ''
  formData.event = 'none'
  formData.images = []
  formData.colors = [{id: undefined, name: ''}]
  imagePreviewUrls.value = []

  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })
}

const closeModals = () => {
  setOpenDialog(false, {name: 'CreationDialog'})
  resetForm()
}

const emit = defineEmits<{
  productAdded: [product: Product]
}>()

const createFormData = (): FormData => {
  const formDataToSend = new FormData()

  const productData = {
    name: formData.name.trim(),
    description: formData.description.trim(),
    price: Number(formData.price).toFixed(2),
    stock: Number(formData.amount),
    model: formData.model.trim().toLowerCase(),
    event: formData.event,
    colors: formData.colors
        .filter(color => color.name.trim())
        .map(color => ({
          id: color.id || undefined,
          name: color.name.trim(),
          hexCode: null
        })),
    metadata: {
      createdAt: new Date().toISOString(),
      createdBy: 'admin',
      version: '1.0'
    }
  }

  formDataToSend.append('product', JSON.stringify(productData))

  formData.images.forEach((file, index) => {
    formDataToSend.append(`pictures`, file)
    formDataToSend.append(`image_${index}_isPrimary`, (index === 0).toString())
  })

  return formDataToSend
}

const handleSubmit = async () => {
  if (!validateForm()) {
    errors.general = 'Please fix the errors above before submitting'
    return
  }

  isSubmitting.value = true
  errors.general = ''

  try {
    const formDataPayload = createFormData()

    const res = await addProduct(formDataPayload)

    if (res === 201) {
      toast.success("Product was created successfully.")
      setOpenDialog(false, {name: 'CreationDialog'})
    }
  } catch (error) {
    console.error('Error creating product:', error)
    errors.general = 'Failed to create product. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}

const isFormValid = computed(() => {
  return formData.name.trim() &&
      formData.description.trim() &&
      formData.model.trim() &&
      formData.price > 0 &&
      formData.amount >= 0 &&
      formData.images.length > 0 &&
      formData.colors.some(color => color.name.trim())
})
</script>

<template>
  <div class="bg-white rounded-t-lg sm:rounded-lg w-full p-6 pb-0 sm:max-w-2xl sm:mx-4 max-h-[90vh] overflow-y-auto">
    <form @submit.prevent="handleSubmit" class="space-y-4 sm:space-y-6">
      <div v-if="errors.general" class="p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg">
        {{ errors.general }}
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
        <div class="sm:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Name <span class="text-red-500">*</span>
          </label>
          <input
              v-model="formData.name"
              type="text"
              :class="[
              'w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent',
              errors.name ? 'border-red-500' : 'border-gray-300'
            ]"
          >
          <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Model <span class="text-red-500">*</span>
          </label>
          <select
              v-model="formData.model"
              style="-webkit-appearance: none;"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="bed">Bed</option>
            <option value="chair">Chair</option>
            <option value="table">Table</option>
            <option value="sofa">Sofa</option>
          </select>
          <p v-if="errors.model" class="mt-1 text-sm text-red-600">{{ errors.model }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Event</label>
          <select
              v-model="formData.event"
              style="-webkit-appearance: none;"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="none">None</option>
            <option value="sale">Sale</option>
            <option value="featured">Featured</option>
            <option value="new">New</option>
            <option value="bestseller">Bestseller</option>
          </select>
          <!--          <GeneralSelect :product-event="formData.event" v-model="formData.event" select-value="Select an event" select-label="Events" :events="events"></GeneralSelect>-->
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Price ($) <span class="text-red-500">*</span>
          </label>
          <input
              v-model="formData.price"
              type="number"
              step="0.01"
              min="0"
              :class="[
              'w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent',
              errors.price ? 'border-red-500' : 'border-gray-300'
            ]"
          >
          <p v-if="errors.price" class="mt-1 text-sm text-red-600">{{ errors.price }}</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Initial Stock <span class="text-red-500">*</span>
          </label>
          <input
              v-model="formData.amount"
              type="number"
              min="0"
              :class="[
              'w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent',
              errors.amount ? 'border-red-500' : 'border-gray-300'
            ]"
          >
          <p v-if="errors.amount" class="mt-1 text-sm text-red-600">{{ errors.amount }}</p>
        </div>
        <div class="sm:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Description <span class="text-red-500">*</span>
          </label>
          <textarea
              v-model="formData.description"
              rows="3"
              :class="[
              'w-full px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent',
              errors.description ? 'border-red-500' : 'border-gray-300'
            ]"
          ></textarea>
          <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
        </div>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Product Images <span class="text-red-500">*</span>
        </label>
        <div v-if="formData.images.length > 0" class="grid grid-cols-2 sm:grid-cols-3 gap-4 mb-4">
          <div v-for="(image, index) in formData.images" :key="index" class="relative group">
            <img
                :src="imagePreviewUrls[index]"
                :alt="`Preview ${index + 1}`"
                class="w-full h-24 sm:h-32 object-cover rounded-lg border border-gray-300"
            >
            <button
                @click="removePictureUrl(index)"
                type="button"
                class="absolute -top-2 -right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
            <div v-if="index === 0" class="absolute bottom-1 left-1 bg-blue-500 text-white text-xs px-2 py-1 rounded">
              Primary
            </div>
          </div>
        </div>
        <button
            @click="addPicture"
            type="button"
            :class="[
            'w-full px-4 py-8 border-2 border-dashed rounded-lg text-sm font-medium transition-colors',
            errors.images
              ? 'border-red-300 text-red-600 hover:bg-red-50'
              : 'border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400'
          ]"
        >
          <svg class="mx-auto h-8 w-8 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          {{ formData.images.length === 0 ? 'Upload Product Images' : 'Add More Images' }}
          <p class="text-xs text-gray-500 mt-1">
            JPEG, PNG, WebP, GIF up to 5MB each
          </p>
        </button>

        <p v-if="errors.images" class="mt-1 text-sm text-red-600">{{ errors.images }}</p>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Colors <span class="text-red-500">*</span>
        </label>
        <div class="flex flex-wrap gap-2 pb-2">
          <div @click="addColor(color)"
               :class="cn('flex cursor-pointer hover:bg-gray-100 items-center justify-center py-1 w-[100px] border rounded-lg', isColorChosen(color.name) && 'border-green-200')"
               :key="index" v-for="(color, index) in allColors">
            <p>{{ color.name }}</p>
          </div>
        </div>
        <div class="space-y-2">
          <div v-for="(color, index) in formData.colors" :key="index" class="flex gap-2">
            <input
                v-model="color.name"
                type="text"
                placeholder="New color name (e.g., Red, Blue, Black)"
                :class="[
                'flex-1 px-3 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm',
                errors.colors ? 'border-red-500' : 'border-gray-300'
              ]"
            >
            <button
                v-if="formData.colors.length > 1"
                @click="removeColor(index)"
                type="button"
                class="px-3 py-2 text-red-600 hover:text-red-800 hover:bg-red-50 rounded-lg"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
          <button
              @click="addColor()"
              type="button"
              class="w-full px-3 py-2 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:bg-gray-50"
          >
            + Add another color
          </button>
        </div>
        <p v-if="errors.colors" class="mt-1 text-sm text-red-600">{{ errors.colors }}</p>
      </div>
      <div class="sticky bottom-0 bg-white pt-4 border-t border-gray-200 flex flex-col sm:flex-row gap-2 sm:gap-3">
        <button
            @click="closeModals"
            type="button"
            :disabled="isSubmitting"
            class="flex-1 px-4 py-3 sm:py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium disabled:opacity-50"
        >
          Cancel
        </button>
        <button
            type="submit"
            :disabled="isSubmitting || !isFormValid"
            class="flex-1 px-4 py-3 sm:py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg"
               fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isSubmitting ? 'Creating...' : 'Add Product' }}
        </button>
      </div>
    </form>
  </div>
</template>