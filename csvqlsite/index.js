document.addEventListener('alpine:init', () => {
    Alpine.data("app", () => ({
        result: null,
        query: "",
        loading: false,
        error: "",
        schemas: null,
        ready: false,
        open: false,


        async init() {
            const response = await fetch('/csvql/schemas')
            this.schemas = await response.json()
            this.ready = true
        },

        async execute() {
            this.error = ""
            this.loading = true
            this.rows = []
            this.columns = []

            try {
                const response = await fetch('/csvql/query', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: this.query })
                })

                if (!response.ok) {
                    const data = await response.json()
                    this.error = data.error || 'An error occurred'
                }

                this.result = await response.json()
            } catch (e) {
                this.error = e.message
            } finally {
                this.loading = false
            }
        }

    }))
})