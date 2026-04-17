import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useChatStore } from '@/stores/chat'
import type { FormFillData } from '@/services/chatService'

// ---------------------------------------------------------------------------
// Shared fixture
// ---------------------------------------------------------------------------
const sampleFill: FormFillData = {
  title: 'Add Payment Method X on website',
  businessLine: '18518',
  requesterBU: 'ADEO-8095',
  busInterested: ['20059', '20060'],
  timeSensitive: 'No',
  whyDemand: 'Current payment options are limited',
  whoIsImpacted: 'All end customers, ~2M users',
  benefitCategory: 'Customer satisfaction & revenue',
  benefitHypothesis: 'Adding this will increase conversion by 15%',
  measureBenefits: 'Track conversion rate monthly',
}

const secondFill: FormFillData = {
  title: 'Improve Search Results',
  businessLine: '18519',
  requesterBU: 'ADEO-8064',
  busInterested: ['20047'],
  timeSensitive: 'Legal',
  whyDemand: 'Search is slow and inaccurate',
  whoIsImpacted: 'Internal operators, ~500 users',
  benefitCategory: 'Service quality & security risk',
  benefitHypothesis: 'Better indexing reduces search time by 40%',
  measureBenefits: 'Average search response time and NPS',
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
describe('useChatStore', () => {
  beforeEach(() => {
    // Create a fresh Pinia instance before each test to avoid state leakage
    setActivePinia(createPinia())
  })

  describe('initial state', () => {
    it('starts with pendingFormFill as null', () => {
      const store = useChatStore()
      expect(store.pendingFormFill).toBeNull()
    })
  })

  describe('setFormFill', () => {
    it('sets pendingFormFill to the provided data', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)
      expect(store.pendingFormFill).toEqual(sampleFill)
    })

    it('stores a reference-equal copy of the provided object', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)
      // pendingFormFill should deeply equal the input
      expect(store.pendingFormFill?.title).toBe(sampleFill.title)
      expect(store.pendingFormFill?.busInterested).toEqual(sampleFill.busInterested)
    })

    it('overwrites a previous value when called again', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)
      expect(store.pendingFormFill?.title).toBe('Add Payment Method X on website')

      store.setFormFill(secondFill)
      expect(store.pendingFormFill?.title).toBe('Improve Search Results')
      expect(store.pendingFormFill?.businessLine).toBe('18519')
    })

    it('handles multiple rapid successive calls, preserving only the last', () => {
      const store = useChatStore()

      store.setFormFill(sampleFill)
      store.setFormFill(secondFill)
      store.setFormFill(sampleFill)

      // Last call wins
      expect(store.pendingFormFill).toEqual(sampleFill)
    })

    it('handles a fill with empty busInterested array', () => {
      const store = useChatStore()
      const fillWithEmptyBus: FormFillData = { ...sampleFill, busInterested: [] }
      store.setFormFill(fillWithEmptyBus)
      expect(store.pendingFormFill?.busInterested).toEqual([])
    })

    it('stores all required fields correctly', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)

      expect(store.pendingFormFill?.title).toBe(sampleFill.title)
      expect(store.pendingFormFill?.businessLine).toBe(sampleFill.businessLine)
      expect(store.pendingFormFill?.requesterBU).toBe(sampleFill.requesterBU)
      expect(store.pendingFormFill?.busInterested).toEqual(sampleFill.busInterested)
      expect(store.pendingFormFill?.timeSensitive).toBe(sampleFill.timeSensitive)
      expect(store.pendingFormFill?.whyDemand).toBe(sampleFill.whyDemand)
      expect(store.pendingFormFill?.whoIsImpacted).toBe(sampleFill.whoIsImpacted)
      expect(store.pendingFormFill?.benefitCategory).toBe(sampleFill.benefitCategory)
      expect(store.pendingFormFill?.benefitHypothesis).toBe(sampleFill.benefitHypothesis)
      expect(store.pendingFormFill?.measureBenefits).toBe(sampleFill.measureBenefits)
    })
  })

  describe('clearFormFill', () => {
    it('sets pendingFormFill back to null', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)
      expect(store.pendingFormFill).not.toBeNull()

      store.clearFormFill()
      expect(store.pendingFormFill).toBeNull()
    })

    it('is a no-op when called on an already-null state', () => {
      const store = useChatStore()
      // pendingFormFill is already null; calling clear should not throw
      expect(() => store.clearFormFill()).not.toThrow()
      expect(store.pendingFormFill).toBeNull()
    })

    it('can be called multiple times without error', () => {
      const store = useChatStore()
      store.setFormFill(sampleFill)
      store.clearFormFill()
      store.clearFormFill()
      store.clearFormFill()
      expect(store.pendingFormFill).toBeNull()
    })
  })

  describe('setFormFill → clearFormFill lifecycle', () => {
    it('set then clear then set again works correctly', () => {
      const store = useChatStore()

      store.setFormFill(sampleFill)
      expect(store.pendingFormFill).not.toBeNull()

      store.clearFormFill()
      expect(store.pendingFormFill).toBeNull()

      store.setFormFill(secondFill)
      expect(store.pendingFormFill?.title).toBe('Improve Search Results')
    })

    it('two separate store instances share state via Pinia singleton', () => {
      const store1 = useChatStore()
      const store2 = useChatStore()

      store1.setFormFill(sampleFill)
      // Both references point to the same store
      expect(store2.pendingFormFill?.title).toBe(sampleFill.title)

      store2.clearFormFill()
      expect(store1.pendingFormFill).toBeNull()
    })
  })
})
