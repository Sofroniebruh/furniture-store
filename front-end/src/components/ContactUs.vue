<script setup>
import {Button} from "@/components/ui/button/index.js";
import {z} from "zod";
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import Wrapper from "@/components/Wrapper.vue";

const schema = z.object({
  firstName: z.string().min(1, {message: 'First name is required'}),
  lastName: z.string().min(1, {message: 'Last name is required'}),
  email: z.string().min(1, {message: "Email is required"}).email({message: 'Email is not valid'}),
  message: z.string().min(1, {message: 'Message is required'}),
})

const {handleSubmit, errors, resetForm} = useForm({
  validationSchema: toTypedSchema(schema),
  validateOnMount: false,
  initialValues: {
    firstName: '',
    lastName: '',
    email: '',
    message: ''
  }
})

const {value: firstName} = useField('firstName')
const {value: lastName} = useField('lastName')
const {value: email} = useField('email')
const {value: message} = useField('message')

const onSubmit = handleSubmit((values) => {
  console.log('Form submitted:', values)
  resetForm()
})
</script>

<template>
  <div class="w-full min-h-[calc(100vh-68px-200px)] flex items-center justify-center flex-col">
    <h1 class="sm:text-4xl text-3xl sm:mt-[68px] mt-[34px]">Contact Us</h1>
    <Wrapper class="w-full sm:w-[400px] mt-[34px] sm:mt-[68px]">
      <form class="flex flex-col gap-2 w-full" @submit="onSubmit">
        <label for="firstName">First name</label>
        <div class="w-full flex flex-col">
          <input
              id="firstName"
              class="border rounded-sm px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
              v-model="firstName"
              type="text"
              name="firstName"
              placeholder="John"
          />
          <p v-if="errors.firstName" class="text-red-600 text-sm mt-1">{{ errors.firstName }}</p>
        </div>
        <label for="lastName">Last name</label>
        <div class="w-full flex flex-col">
          <input
              id="lastName"
              class="border rounded-sm px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
              v-model="lastName"
              type="text"
              name="lastName"
              placeholder="Doe"
          />
          <p v-if="errors.lastName" class="text-red-600 text-sm mt-1">{{ errors.lastName }}</p>
        </div>
        <label for="email">Email</label>
        <div class="w-full flex flex-col">
          <input
              id="email"
              class="border rounded-sm px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
              v-model="email"
              type="email"
              name="email"
              placeholder="example@mail.com"
          />
          <p v-if="errors.email" class="text-red-600 text-sm mt-1">{{ errors.email }}</p>
        </div>
        <label for="message">Message</label>
        <div class="w-full flex flex-col">
          <textarea
              id="message"
              class="border rounded-sm px-2 py-1 min-h-[130px] focus:outline-none focus:ring-2 focus:ring-blue-500"
              :rows="5"
              v-model="message"
              name="message"
              placeholder="Your message..."
          />
          <p v-if="errors.message" class="text-red-600 text-sm mt-1">{{ errors.message }}</p>
        </div>
        <Button class="cursor-pointer mt-[17px]" type="submit">Send</Button>
      </form>
    </Wrapper>
  </div>
</template>