#include <dlfcn.h>
#include <stdlib.h>

#include "linux_cups.h"

/*
 * dynamic implementation of some cups functions
 * Assumes functionality from LSB Printing 3.2 (2007)
 * http://refspecs.linuxfoundation.org/LSB_3.2.0/LSB-Printing/LSB-Printing/book1.html
 *
 * This allows things to be compiled without the CUPS headers installed, just requiring -ldl
 */

static void *cupslib = NULL;

static int (*cupsgetdests)(cups_dest_t **);
static void (*cupsfreedests)(int, cups_dest_t *);
static cups_dest_t *(*cupsgetdest)(const char *, const char *, int, cups_dest_t *);
static int (*cupsaddoption)(const char *, const char *, int, cups_option_t **);
static void (*cupsfreeoptions)(int, cups_option_t *);
static int (*cupsprintfile)(const char *, const char *, const char *, int, cups_option_t *);

int linux_cups_init_dl(void) {
#define fail_if_null(x) if (x == NULL) { return 0; }
	cupslib = dlopen("libcups.so.2", RTLD_LAZY);
	fail_if_null(cupslib);
	cupsgetdests = dlsym(cupslib, "cupsGetDests");
	fail_if_null(cupsgetdests)
	cupsfreedests = dlsym(cupslib, "cupsFreeDests");
	fail_if_null(cupsfreedests)
	cupsgetdest = dlsym(cupslib, "cupsGetDest");
	fail_if_null(cupsgetdest)
	cupsaddoption = dlsym(cupslib, "cupsAddOption");
	fail_if_null(cupsaddoption)
	cupsfreeoptions = dlsym(cupslib, "cupsFreeOptions");
	fail_if_null(cupsfreeoptions)
	cupsprintfile = dlsym(cupslib, "cupsPrintFile");
	fail_if_null(cupsprintfile)
#undef fail_if_null
	return 1;
}

int linux_cups_get_dests(cups_dest_t **dests) {
	return cupsgetdests(dests);
}

void linux_cups_free_dests(int num_dests, cups_dest_t *dests) {
	cupsfreedests(num_dests, dests);
}

cups_dest_t *linux_cups_get_dest(const char *name, const char *instance, int num_dests, cups_dest_t *dests) {
	return cupsgetdest(name, instance, num_dests, dests);
}

int linux_cups_add_option(const char *name, const char *value, int num_options, cups_option_t **options) {
	return cupsaddoption(name, value, num_options, options);
}

void linux_cups_free_options(int num_opts, cups_option_t *opts) {
	cupsfreeoptions(num_opts, opts);
}

int linux_cups_print_file(const char *name, const char *filename, const char *title, int num_options, cups_option_t *options) {
	return cupsprintfile(name, filename, title, num_options, options);
}
