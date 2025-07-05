// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  integrations: [
    starlight({
      title: "Keezle",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/gaurishhs/keezle",
        },
      ],
      sidebar: [
        {
          label: "Getting Started",
          slug: "getting-started",
        },
        {
          label: "Reference",
          autogenerate: { directory: "reference" },
        },
      ],
    }),
  ],
});
