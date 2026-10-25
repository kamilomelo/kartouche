# Updating a copied kartouche template

For a coding agent updating a host script or repository:

1. Locate its `kartouche.<ext>` file, or the pasted region between `kartouche:code begin vN` and `kartouche:code end`.
2. Fetch the latest template from a local clone of [kartouche](https://github.com/kamilomelo/kartouche), or from `https://raw.githubusercontent.com/kamilomelo/kartouche/main/<language>/kartouche.<ext>` (for example, [the Bash template](https://raw.githubusercontent.com/kamilomelo/kartouche/main/bash/kartouche.sh)).
3. Compare the `vN` in the `kartouche:code begin` markers. If equal, stop and report “already up to date”.
4. Replace the **whole code region including both markers** with the latest code region. Update only the header's `KARTOUCHE_TEMPLATE_VERSION` value to match the new marker.
5. Preserve the `kartouche:config` region, the host's code, and everything else outside the code region, except the header version. For Go only, reconcile the `import` block before config: add or remove imports needed by the new template code, but keep any imports used elsewhere in the host.
6. If an existing host call uses a removed or renamed function, identify and fix every affected call site. If the new template needs additional config fields, add their defaults to the host's config and tell the owner. Don't overwrite customized values.
7. Run the host's tests (or at least a dry run), show the diff, and report any remaining issues.

The copy belongs to the host project; do not turn it into a shared dependency.
