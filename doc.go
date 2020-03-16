// Licensed under the GPLv3, see LICENCE file for details.

/*
This package provides a simple and effective way to collect errors while
simplifies the flow complexity of sequential and conditional safety checks.

The exported 'NewErr' and 'NewErrOr' alows initialization with or without a
previous 'Errs'. Theses functions work like a wrap around the built-in and
'juju/errors' constructor functions.

Once initialized, 'errs.NewErr' and 'errs.NewErrWithCause' work as a
replacement of 'errors.NewErr' and 'errors.NewErrWithCause' that appends the
generated jujuErr to the inner errors. Also the function 'errs.Append' appends
the submitted error arguments to the inner slice.

A primary use case for this library is to append multiple errors while
doing a sequence of checkings.

  errs := make([]error, 0)
  if err := checkFunc1(); err != nil {
    errs = append(errs, err)
  }

  if err := checkFunc2(); err != nil {
    errs = append(errs, err)
  }

Would become with github.com/juju/errors:

  errs := make([]error, 0)
  if err := checkFunc1(); err != nil {
    errs = append(errs, errors.Trace(err))
  }

  if err := checkFunc2(); err != nil {
    errs = append(errs, errors.Annotate(err, "more context"))
    // Adding annotation to the error
  }

And with bulkerrs:

  errs := bulkerrs.NewErrOr(checkFunc1())
  errs.NewErrWithCause(checkFunc2(), "more context")

There's no longer need to check if the error is nil.

Additionally, bulkerrs makes easy to integrate the errors appendings in the
application control flow:

  errs := make([]error, 0)
  if err1 := checkFunc1(); err1 != nil {
    errs = append(errs, err1)
    if err1_1 := checkFunc1_1(); err1_1 != nil {
      errs = append(errs, err1_1)
    }
  } else if err2 := checkFunc2(); err2 != nil {
    errs = append(errs, err2)
    if err2_1 := checkFunc2_1(); err2_1 != nil {
      errs = append(errs, err2_1)
    }
  } else {
    if errx_1 := checkFuncx(); errx_1 != nil {
      errs = append(errs, errx_1)
    } else {
      return nil
    }
  }
  return &errs

Would become:

  errs := bulkerrs.NewErr()
  if errs.Append(checkFunc1()) {
    errs.Append(checkFunc1_1())
  } else if errs.Append(checkFunc2()) {
    errs.Append(checkFunc2_1())
  } else {
    errs.Append(checkFuncx())
  }
  // return an error if len(errs.errors) > 0 or nil
  return errs.ToError()

And if needed, like in github.com/juju/errors, it's possible to add extra
context, and have an advanced control of the application flow:

  errs := bulkerrs.NewErr()
  if err1 := checkFunc1(); errs.AppendIfX(err1 != nil, errors.Annotate, "more context1", err1) {
    errs.Append(checkFunc1_1())
  } else if err2 := checkFunc2(); errs.AppendIfX(err2 == nil, errors.Annotate, "This should have failed but didn't", nil) {
    errs.Append(checkFunc2_1())
  } else {
    errx := checkFuncx()
    // add extra error types
    errs.AppendIfX(errx != nil, errors.NewNotValid, "more context", errx)
  }
  return errs.ToError()


When you want to check to see if an error is of a particular type, a helper
function is normally exported by the package that returned the error, like the
'os' package does.  The underlying cause of the error is available using the
'Cause' function.
	os.IsNotExist(errors.Cause(err))
The result of the 'Error()' call on an annotated error is the annotations joined
with colons, then the result of the 'Error()' method for the underlying error
that was the cause.
	err := errors.Errorf("original")
	err = errors.Annotatef(err, "context")
	err = errors.Annotatef(err, "more context")
	err.Error() -> "more context: context: original"
Obviously recording the file, line and functions is not very useful if you
cannot get them back out again.
	errors.ErrorStack(err)
will return something like:
	first error
	github.com/juju/errors/annotation_test.go:193:
	github.com/juju/errors/annotation_test.go:194: annotation
	github.com/juju/errors/annotation_test.go:195:
	github.com/juju/errors/annotation_test.go:196: more context
	github.com/juju/errors/annotation_test.go:197:
The first error was generated by an external system, so there was no location
associated. The second, fourth, and last lines were generated with Trace calls,
and the other two through Annotate.
Sometimes when responding to an error you want to return a more specific error
for the situation.
    if err := FindField(field); err != nil {
	    return errors.Wrap(err, errors.NotFoundf(field))
	}
This returns an error where the complete error stack is still available, and
'errors.Cause()' will return the 'NotFound' error.

*/
package bulkerrs
