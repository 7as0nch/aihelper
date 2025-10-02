<template>
  <form class="form" @submit.prevent="handleSubmit">
    <div v-for="field in fields" :key="field.name" class="form-group">
      <label :for="field.name" class="form-label">{{ field.label }}</label>
      <input
        v-if="field.type === 'text' || field.type === 'email' || field.type === 'password'"
        :id="field.name"
        :type="field.type"
        :name="field.name"
        v-model="formData[field.name]"
        :required="field.required"
        class="form-input"
      />
      <textarea
        v-else-if="field.type === 'textarea'"
        :id="field.name"
        :name="field.name"
        v-model="formData[field.name]"
        :required="field.required"
        class="form-textarea"
      ></textarea>
      <select
        v-else-if="field.type === 'select'"
        :id="field.name"
        :name="field.name"
        v-model="formData[field.name]"
        :required="field.required"
        class="form-select"
      >
        <option v-for="option in field.options" :key="option.value" :value="option.value">
          {{ option.label }}
        </option>
      </select>
    </div>
    <div class="form-actions">
      <button type="submit" class="btn-submit">{{ submitText }}</button>
      <button type="button" class="btn-cancel" @click="handleCancel">{{ cancelText }}</button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { reactive, defineProps, defineEmits } from 'vue';

interface Field {
  name: string;
  label: string;
  type: 'text' | 'email' | 'password' | 'textarea' | 'select';
  required?: boolean;
  options?: { value: string; label: string }[];
}

interface FormData {
  [key: string]: string;
}

const props = defineProps<{
  fields: Field[];
  submitText?: string;
  cancelText?: string;
}>();

const emit = defineEmits<{
  (e: 'submit', data: FormData): void;
  (e: 'cancel'): void;
}>();

const formData = reactive<FormData>({});

// 初始化表单数据
props.fields.forEach(field => {
  formData[field.name] = '';
});

const handleSubmit = () => {
  emit('submit', formData);
};

const handleCancel = () => {
  emit('cancel');
};
</script>

<style scoped>
.form {
  background-color: #ffffff;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #495057;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #42b983;
  box-shadow: 0 0 0 3px rgba(66, 185, 131, 0.2);
}

.form-textarea {
  min-height: 120px;
  resize: vertical;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.btn-submit,
.btn-cancel {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-submit {
  background-color: #42b983;
  color: white;
}

.btn-submit:hover {
  background-color: #359c6d;
}

.btn-cancel {
  background-color: #6c757d;
  color: white;
}

.btn-cancel:hover {
  background-color: #5a6268;
}
</style>