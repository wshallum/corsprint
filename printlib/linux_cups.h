#ifndef PRINTLIB_LINUX_CUPS_H_INCLUDED
#define PRINTLIB_LINUX_CUPS_H_INCLUDED
typedef struct {
  char *name;
  char *value;
} cups_option_t;
typedef struct {
  char *name;
  char *instance;
  int is_default;
  int num_options;
  cups_option_t *options;
} cups_dest_t;
int linux_cups_init_dl(void);
int linux_cups_get_dests(cups_dest_t **);
void linux_cups_free_dests(int, cups_dest_t *);
cups_dest_t *linux_cups_get_dest(const char *, const char *, int, cups_dest_t *);
int linux_cups_add_option(const char *, const char *, int, cups_option_t **);
void linux_cups_free_options(int, cups_option_t *);
int linux_cups_print_file(const char *, const char *, const char *, int, cups_option_t *);
#endif
