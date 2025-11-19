import { useEffect, useMemo, useRef, useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'

type ManualFormState = {
  destinations: string
  subject: string
  body: string
  template: string
}

type CsvFormState = {
  subject: string
  body: string
  template: string
  targetColumn: string
  file: File | null
}

type StatusMessage = {
  type: 'success' | 'error'
  message: string
}

const API_BASE_URL = (() => {
  const raw = import.meta.env.VITE_API_BASE_URL ?? ''
  if (!raw) return ''
  return raw.replace(/\/+$/, '')
})()

const TemplatePreview = ({ label, value }: { label: string; value: string }) => {
  const trimmed = value.trim()
  const iframeRef = useRef<HTMLIFrameElement | null>(null)
  const [frameHeight, setFrameHeight] = useState(360)

  const srcDoc = useMemo(() => {
    if (!trimmed) return ''
    return `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <style>
      body {
        font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
        padding: 24px;
        color: #0f0f0f;
        background: #ffffff;
        line-height: 1.6;
        margin: 0;
      }
    </style>
  </head>
  <body>${trimmed}</body>
</html>`
  }, [trimmed])

  useEffect(() => {
    const iframe = iframeRef.current
    if (!iframe) return
    let resizeObserver: ResizeObserver | null = null

    const syncHeight = () => {
      const doc = iframe.contentDocument
      const height =
        doc?.body?.scrollHeight || doc?.documentElement?.scrollHeight || iframeRef.current?.scrollHeight || 360
      setFrameHeight(Math.max(height, 320))
    }

    const handleLoad = () => {
      syncHeight()
      const target = iframe.contentDocument?.body
      if (target && 'ResizeObserver' in window) {
        resizeObserver = new ResizeObserver(() => syncHeight())
        resizeObserver.observe(target)
      } else {
        // fall back to delayed resize for environments without ResizeObserver support
        setTimeout(() => syncHeight(), 250)
        setTimeout(() => syncHeight(), 750)
      }
    }

    iframe.addEventListener('load', handleLoad)
    return () => {
      iframe.removeEventListener('load', handleLoad)
      resizeObserver?.disconnect()
    }
  }, [srcDoc])

  if (!trimmed) return null

  return (
    <div className="overflow-hidden rounded-3xl border border-white/10 bg-white/5 p-4 shadow-inner shadow-black/40 backdrop-blur">
      <div className="flex items-center justify-between text-[0.65rem] font-semibold uppercase tracking-[0.35em] text-white/60">
        <span>{label}</span>
        <span>live preview</span>
      </div>
      <iframe
        ref={iframeRef}
        title={`${label} iframe`}
        sandbox="allow-same-origin"
        scrolling="no"
        srcDoc={srcDoc}
        className="w-full rounded-2xl border border-white/10 bg-white"
        style={{ height: frameHeight, overflow: 'hidden' }}
      />
    </div>
  )
}

const StatusNotice = ({ status }: { status: StatusMessage | null }) => {
  if (!status) return null
  return (
    <p
      className={`text-sm font-medium ${
        status.type === 'success' ? 'text-emerald-200' : 'text-rose-200'
      }`}
    >
      {status.message}
    </p>
  )
}

const initialManualForm: ManualFormState = {
  destinations: '',
  subject: '',
  body: '',
  template: '',
}

const initialCsvForm: CsvFormState = {
  subject: '',
  body: '',
  template: '',
  targetColumn: '',
  file: null,
}

function App() {
  const [manualForm, setManualForm] = useState<ManualFormState>(initialManualForm)
  const [manualStatus, setManualStatus] = useState<StatusMessage | null>(null)
  const [manualLoading, setManualLoading] = useState(false)

  const [csvForm, setCsvForm] = useState<CsvFormState>(initialCsvForm)
  const [csvStatus, setCsvStatus] = useState<StatusMessage | null>(null)
  const [csvLoading, setCsvLoading] = useState(false)
  const [csvFileInputKey, setCsvFileInputKey] = useState(0)

  const handleManualFieldChange =
    (field: keyof ManualFormState) =>
    (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const { value } = event.target
      setManualForm((prev) => ({
        ...prev,
        [field]: value,
      }))
    }

  const handleCsvFieldChange =
    (field: keyof CsvFormState) =>
    (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
      const { value } = event.target
      setCsvForm((prev) => ({
        ...prev,
        [field]: value,
      }))
    }

  const handleCsvFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null
    setCsvForm((prev) => ({
      ...prev,
      file,
    }))
  }

  const postJson = async (path: string, payload: unknown) => {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      const message =
        typeof (data as { error?: string })?.error === 'string'
          ? (data as { error: string }).error
          : 'Unable to submit request. Please try again.'
      throw new Error(message)
    }
    return data
  }

  const postFormData = async (path: string, payload: FormData) => {
    const response = await fetch(`${API_BASE_URL}${path}`, {
      method: 'POST',
      body: payload,
    })
    const data = await response.json().catch(() => ({}))
    if (!response.ok) {
      const message =
        typeof (data as { error?: string })?.error === 'string'
          ? (data as { error: string }).error
          : 'Unable to submit request. Please try again.'
      throw new Error(message)
    }
    return data
  }

  const handleManualSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setManualStatus(null)

    const destinations = manualForm.destinations
      .split(/[\n,]+/)
      .map((address) => address.trim())
      .filter(Boolean)

    if (!destinations.length) {
      setManualStatus({
        type: 'error',
        message: 'Add at least one destination email.',
      })
      return
    }

    if (!manualForm.subject.trim()) {
      setManualStatus({
        type: 'error',
        message: 'Subject is required.',
      })
      return
    }

    setManualLoading(true)
    try {
      const payload = {
        destinations,
        subject: manualForm.subject.trim(),
        body: manualForm.body.trim(),
        template: manualForm.template.trim(),
      }

      const result = (await postJson('/email', payload)) as { data?: { batchId?: string } }
      const batchId = result?.data?.batchId ?? 'N/A'

      setManualStatus({
        type: 'success',
        message: `Email batch queued successfully (Batch ID: ${batchId}).`,
      })
    } catch (error) {
      setManualStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'Unexpected error.',
      })
    } finally {
      setManualLoading(false)
    }
  }

  const handleCsvSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setCsvStatus(null)

    if (!csvForm.file) {
      setCsvStatus({
        type: 'error',
        message: 'Please attach a CSV file.',
      })
      return
    }

    if (!csvForm.subject.trim() || !csvForm.targetColumn.trim()) {
      setCsvStatus({
        type: 'error',
        message: 'Subject and target column are required.',
      })
      return
    }

    setCsvLoading(true)
    try {
      const formData = new FormData()
      formData.append('subject', csvForm.subject.trim())
      formData.append('body', csvForm.body.trim())
      formData.append('template', csvForm.template.trim())
      formData.append('targetColumn', csvForm.targetColumn.trim())
      formData.append('data', csvForm.file)

      const result = (await postFormData('/email/csv', formData)) as { data?: { batchId?: string } }
      const batchId = result?.data?.batchId ?? 'N/A'

      setCsvStatus({
        type: 'success',
        message: `CSV batch queued successfully (Batch ID: ${batchId}).`,
      })
      setCsvForm((prev) => ({
        ...prev,
        file: null,
      }))
      setCsvFileInputKey((key) => key + 1)
    } catch (error) {
      setCsvStatus({
        type: 'error',
        message: error instanceof Error ? error.message : 'Unexpected error.',
      })
    } finally {
      setCsvLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-[#050607] via-[#0d1015] to-[#151820] px-4 py-10 text-white sm:px-6 lg:px-12">
      <div className="flex w-full flex-col gap-12">
        <header className="space-y-3 text-left">
          <p className="text-[0.7rem] uppercase tracking-[0.6em] text-white/60">
            Transactional Operations Assistant
          </p>
          <h1 className="text-3xl font-semibold leading-tight text-white md:text-5xl">
            Send email campaigns with confidence
          </h1>
          <p className="max-w-2xl text-base text-white/70">
            Use either the dedicated form for a handful of contacts or upload a CSV to broadcast to
            thousands. Everything stays monochrome, focused, and ready for the API living in <code>/api</code>.
          </p>
        </header>

        <main className="grid gap-8 lg:grid-cols-2 lg:items-start">
          <article className="flex flex-col gap-6 rounded-[28px] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/40 backdrop-blur">
            <div className="space-y-3">
              <p className="text-xs uppercase tracking-[0.5em] text-white/60">Manual composer</p>
              <h2 className="text-2xl font-semibold text-white">Targeted send</h2>
              <p className="text-sm text-white/70">
                Ideal for quick runs or QA verifications. Drop comma or newline separated email addresses and
                preview your HTML template before scheduling.
              </p>
            </div>

            <form className="flex flex-col gap-5" onSubmit={handleManualSubmit}>
              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Destinations *</span>
                <textarea
                  required
                  rows={3}
                  placeholder="name@example.com, another@example.com"
                  value={manualForm.destinations}
                  onChange={handleManualFieldChange('destinations')}
                  className="min-h-[3rem] rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
                <span className="text-xs text-white/50">Comma or newline separated list.</span>
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Subject *</span>
                <input
                  required
                  type="text"
                  placeholder="Quarterly product updates"
                  value={manualForm.subject}
                  onChange={handleManualFieldChange('subject')}
                  className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Body</span>
                <textarea
                  rows={3}
                  placeholder="Plain text body if you are not using a template."
                  value={manualForm.body}
                  onChange={handleManualFieldChange('body')}
                  className="min-h-[3rem] rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Template</span>
                <textarea
                  rows={6}
                  placeholder="<h1>Hello there</h1>"
                  value={manualForm.template}
                  onChange={handleManualFieldChange('template')}
                  className="min-h-[6rem] rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <button
                  type="submit"
                  disabled={manualLoading}
                  className="inline-flex items-center justify-center rounded-full bg-white px-6 py-3 text-xs font-semibold uppercase tracking-[0.35em] text-gray-900 transition hover:-translate-y-0.5 disabled:opacity-50 disabled:hover:translate-y-0"
                >
                  {manualLoading ? 'Scheduling…' : 'Send email batch'}
                </button>
                <StatusNotice status={manualStatus} />
              </div>
            </form>

            <TemplatePreview label="Template preview" value={manualForm.template} />
          </article>

          <article className="flex flex-col gap-6 rounded-[28px] border border-white/10 bg-white/5 p-6 shadow-2xl shadow-black/40 backdrop-blur">
            <div className="space-y-3">
              <p className="text-xs uppercase tracking-[0.5em] text-white/60">CSV automation</p>
              <h2 className="text-2xl font-semibold text-white">Bulk upload</h2>
              <p className="text-sm text-white/70">
                Upload a comma separated file, tell us which column holds the target emails, and the API will
                handle the rest. Template preview keeps rendering in lockstep.
              </p>
            </div>

            <form className="flex flex-col gap-5" onSubmit={handleCsvSubmit}>
              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>CSV file *</span>
                <input
                  key={csvFileInputKey}
                  required
                  type="file"
                  accept=".csv,text/csv"
                  onChange={handleCsvFileChange}
                  className="block w-full cursor-pointer rounded-2xl border border-white/10 bg-transparent px-4 py-3 text-sm text-white file:mr-4 file:rounded-full file:border-0 file:bg-white file:px-4 file:py-2 file:text-xs file:font-semibold file:uppercase file:tracking-[0.35em] file:text-gray-900 focus:border-white/40 focus:outline-none"
                />
                <span className="text-xs text-white/50">Headers are required. Attachments stay in-memory only.</span>
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Target column *</span>
                <input
                  required
                  type="text"
                  placeholder="email"
                  value={csvForm.targetColumn}
                  onChange={handleCsvFieldChange('targetColumn')}
                  className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
                <span className="text-xs text-white/50">The column that contains the destination emails.</span>
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Subject *</span>
                <input
                  required
                  type="text"
                  placeholder="New feature rollout"
                  value={csvForm.subject}
                  onChange={handleCsvFieldChange('subject')}
                  className="rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Body</span>
                <textarea
                  rows={3}
                  placeholder="Fallback plain text body."
                  value={csvForm.body}
                  onChange={handleCsvFieldChange('body')}
                  className="min-h-[3rem] rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <label className="flex flex-col gap-2 text-sm font-medium text-white">
                <span>Template</span>
                <textarea
                  rows={6}
                  placeholder="<table role='presentation'>…</table>"
                  value={csvForm.template}
                  onChange={handleCsvFieldChange('template')}
                  className="min-h-[6rem] rounded-2xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder:text-white/40 focus:border-white/40 focus:outline-none"
                />
              </label>

              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <button
                  type="submit"
                  disabled={csvLoading}
                  className="inline-flex items-center justify-center rounded-full bg-white px-6 py-3 text-xs font-semibold uppercase tracking-[0.35em] text-gray-900 transition hover:-translate-y-0.5 disabled:opacity-50 disabled:hover:translate-y-0"
                >
                  {csvLoading ? 'Uploading…' : 'Queue CSV batch'}
                </button>
                <StatusNotice status={csvStatus} />
              </div>
            </form>

            <TemplatePreview label="Template preview" value={csvForm.template} />
          </article>
        </main>
      </div>
    </div>
  )
}

export default App
