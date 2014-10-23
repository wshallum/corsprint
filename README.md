corsprint
=========

Instead of running a Java Applet (with code signing requirements, frequent
security updates, etc.) or an NPAPI plugin (on its way out) to print raw, can
we make requests to a local daemon instead? Modern browsers do CORS, this
should be doable?

Note: Author's first Go program. Likely to be full of stupid decisions. Use at
own risk.

### Other approaches

* Java / ActiveX applet
* Browser extension
* URL protocol handler


### Raw print helpers that actually work

* qz-print : http://qzindustries.com/ (uses Java applet)
* Other implementations are available
