<script setup lang="ts">
import { ref } from 'vue'
import { useAdmin, type Target } from '../../composables/useAdmin'

const emit = defineEmits<{ (e: 'new'): void; (e: 'edit', t: Target): void }>()
const { targets, grants, duplicateTarget, deleteTarget } = useAdmin()

const inlineEditId = ref<number | null>(null)
const inlineEditField = ref<'name' | 'path' | null>(null)

function countGrantsForTarget(id: number) {
  return grants.value.filter(g => g.admin || g.allTargets || g.targetIds.includes(id)).length
}

function startEdit(id: number, field: 'name' | 'path') {
  inlineEditId.value = id
  inlineEditField.value = field
}

function finishEdit() {
  inlineEditId.value = null
  inlineEditField.value = null
}

function commitInline(t: Target, e: Event) {
  const value = (e.target as HTMLInputElement).value.trim()
  if (!value) {
    finishEdit()
    return
  }
  if (inlineEditField.value === 'name' && value !== t.name) {
    emit('edit', { ...t, name: value })
  } else if (inlineEditField.value === 'path' && value !== t.path) {
    emit('edit', { ...t, path: value })
  }
  finishEdit()
}
</script>

<template>
  <div class="fb-fade">
    <div class="section-head">
      <div>
        <h1>Upload targets</h1>
        <div class="sub">Destinations users can upload to. Each target maps a friendly name to a folder path on the storage backend. Click the name or path to rename inline.</div>
      </div>
      <button class="btn btn-primary" @click="emit('new')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
        New target
      </button>
    </div>

    <div v-if="targets.length === 0" class="empty">
      No targets yet. Add one to let people start uploading.
    </div>

    <div v-else class="card">
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Folder path</th>
            <th>Access</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in targets" :key="t.id">
            <td>
              <div class="name-cell">
                <div class="swatch">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/></svg>
                </div>
                <div>
                  <input
                    v-if="inlineEditId === t.id && inlineEditField === 'name'"
                    class="inline-edit"
                    :value="t.name"
                    @blur="commitInline(t, $event)"
                    @keyup.enter="commitInline(t, $event)"
                    @keyup.escape="finishEdit"
                    ref="inlineInput"
                    autofocus
                  />
                  <div v-else class="primary" @click="startEdit(t.id, 'name')" style="cursor:text">{{ t.name }}</div>
                  <div class="secondary">
                    {{ countGrantsForTarget(t.id) }} {{ countGrantsForTarget(t.id) === 1 ? 'grant' : 'grants' }}
                  </div>
                </div>
              </div>
            </td>
            <td>
              <input
                v-if="inlineEditId === t.id && inlineEditField === 'path'"
                class="inline-edit mono"
                :value="t.path"
                @blur="commitInline(t, $event)"
                @keyup.enter="commitInline(t, $event)"
                @keyup.escape="finishEdit"
                autofocus
              />
              <span v-else class="path" @click="startEdit(t.id, 'path')" style="cursor:text">{{ t.path }}</span>
            </td>
            <td>
              <span class="chip" v-if="countGrantsForTarget(t.id) === 0" style="color:var(--ink-3)">No one</span>
              <span v-else class="chip ok">{{ countGrantsForTarget(t.id) }} principals</span>
            </td>
            <td class="actions">
              <button class="btn btn-sm btn-ghost" @click="duplicateTarget(t)">Duplicate</button>
              <button class="btn btn-sm btn-danger" @click="deleteTarget(t.id)">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
