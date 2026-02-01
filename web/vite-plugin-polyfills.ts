import { Plugin } from "vite";

/**
 * Vite plugin to inject core-js polyfills as a traditional (non-module) script
 * that loads before all other scripts to ensure __core-js_shared__ is available.
 */
export function polyfillsPlugin(): Plugin {
  return {
    name: "inject-polyfills",
    transformIndexHtml(html) {
      // Find the polyfills chunk from the generated HTML
      const polyfillsMatch = html.match(/<link rel="modulepreload"[^>]*href="\/assets\/(polyfills-[^"]+)"/);
      if (!polyfillsMatch) return html;

      const polyfillsFile = polyfillsMatch[1];
      // Remove the modulepreload link
      html = html.replace(/<link rel="modulepreload"[^>]*href="\/assets\/polyfills-[^"]+"[^>]*>\n?/, '');
      // Inject traditional script tag for polyfills before module scripts
      const polyfillScript = '<script src="/assets/' + polyfillsFile + '"></script>';
      return html.replace(
        /(<script type="module"[^>]*>)/,
        `${polyfillScript}\n    $1`
      );
    },
  };
}
