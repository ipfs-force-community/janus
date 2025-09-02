# Janus Frontend

The frontend provides a web interface for exploring **Filecoin network upgrades**, including:

* Tracking upgrade history
* Searching and filtering upgrades
* Data visualization

---

## Development

```bash
npm install
npm run dev
````

Then open [http://localhost:3000](http://localhost:3000) in your browser.

---

## Production

Build the project:

```bash
npm run build
```

Start the production server:

```bash
npm run start
```

By default, the app runs on [http://localhost:3000](http://localhost:3000).
To use a custom port:

```bash
PORT=4000 npm run start
```