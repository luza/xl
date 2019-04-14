package formula

import (
	"strings"

	"xl/document/eval"

	"github.com/shopspring/decimal"
)

const maxArguments = 1000

type Function func(*eval.Context, []eval.Value) (eval.Value, error)

type functionDef struct {
	F       Function
	MinArgs int
	MaxArgs int
}

var functions = map[string]functionDef{
	"TRIM": {trim, 1, 1},
	"SUM":  {sum, 1, maxArguments},
	"IF":   {if_, 3, 3},
	// ABS [Math and trigonometry] Returns the absolute value of a number
	// ACCRINT [Financial] Returns the accrued interest for a security that pays periodic interest
	// ACCRINTM [Financial] Returns the accrued interest for a security that pays interest at maturity
	// ACOS [Math and trigonometry] Returns the arccosine of a number
	// ACOSH [Math and trigonometry] Returns the inverse hyperbolic cosine of a number
	// ACOT [Math and trigonometry] Returns the arccotangent of a number
	// ACOTH [Math and trigonometry] Returns the hyperbolic arccotangent of a number
	// AGGREGATE [Math and trigonometry] Returns an aggregate in a list or database
	// ADDRESS [Lookup and reference] Returns a reference as text to a single cell in a worksheet
	// AMORDEGRC [Financial] Returns the depreciation for each accounting period by using a depreciation coefficient
	// AMORLINC [Financial] Returns the depreciation for each accounting period
	// AND [Logical] Returns TRUE if all of its arguments are TRUE
	// ARABIC [Math and trigonometry] Converts a Roman number to Arabic, as a number
	// AREAS [Lookup and reference] Returns the number of areas in a reference
	// ASC [Text] Changes full-width (double-byte) English letters or katakana within a character string to half-width (single-byte) characters
	// ASIN [Math and trigonometry] Returns the arcsine of a number
	// ASINH [Math and trigonometry] Returns the inverse hyperbolic sine of a number
	// ATAN [Math and trigonometry] Returns the arctangent of a number
	// ATAN2 [Math and trigonometry] Returns the arctangent from x- and y-coordinates
	// ATANH [Math and trigonometry] Returns the inverse hyperbolic tangent of a number
	// AVEDEV [Statistical] Returns the average of the absolute deviations of data points from their mean
	// AVERAGE [Statistical] Returns the average of its arguments
	// AVERAGEA [Statistical] Returns the average of its arguments, including numbers, text, and logical values
	// AVERAGEIF [Statistical] Returns the average (arithmetic mean) of all the cells in a range that meet a given criteria
	// AVERAGEIFS [Statistical] Returns the average (arithmetic mean) of all cells that meet multiple criteria.
	// BAHTTEXT [Text] Converts a number to text, using the ß (baht) currency format
	// BASE [Math and trigonometry] Converts a number into a text representation with the given radix (base)
	// BESSELI [Engineering] Returns the modified Bessel function In(x)
	// BESSELJ [Engineering] Returns the Bessel function Jn(x)
	// BESSELK [Engineering] Returns the modified Bessel function Kn(x)
	// BESSELY [Engineering] Returns the Bessel function Yn(x)
	// BETADIST [Compatibility] Returns the beta cumulative distribution
	// BETA.DIST [Statistical] Returns the beta cumulative distribution
	// BETAINV [Compatibility] Returns the inverse of the cumulative distribution function for a specified beta distribution
	// BETA.INV [Statistical] Returns the inverse of the cumulative distribution function for a specified beta distribution
	// BIN2DEC [Engineering] Converts a binary number to decimal
	// BIN2HEX [Engineering] Converts a binary number to hexadecimal
	// BIN2OCT [Engineering] Converts a binary number to octal
	// BINOMDIST [Compatibility] Returns the individual term binomial distribution probability
	// BINOM.DIST [Statistical] Returns the individual term binomial distribution probability
	// BINOM.DIST.RANGE [Statistical] Returns the probability of a trial result using a binomial distribution
	// BINOM.INV [Statistical] Returns the smallest value for which the cumulative binomial distribution is less than or equal to a criterion value
	// BITAND [Engineering] Returns a 'Bitwise And' of two numbers
	// BITLSHIFT [Engineering] Returns a value number shifted left by shift_amount bits
	// BITOR [Engineering] Returns a bitwise OR of 2 numbers
	// BITRSHIFT [Engineering] Returns a value number shifted right by shift_amount bits
	// BITXOR [Engineering] Returns a bitwise 'Exclusive Or' of two numbers
	// CALL [Add-in and Automation] Calls a procedure in a dynamic link library or code resource
	// CEILING [Math and trigonometry] Rounds a number to the nearest integer or to the nearest multiple of significance
	// CEILING.MATH [Math and trigonometry] Rounds a number up, to the nearest integer or to the nearest multiple of significance
	// CEILING.PRECISE [Math and trigonometry] Rounds a number the nearest integer or to the nearest multiple of significance. Regardless of the sign of the number, the number is rounded up.
	// CELL [Information] Returns information about the formatting, location, or contents of a cell
	// CHAR [Text] Returns the character specified by the code number
	// CHIDIST [Compatibility] Returns the one-tailed probability of the chi-squared distribution
	// Note: CHIINV [Compatibility] Returns the inverse of the one-tailed probability of the chi-squared distribution
	// Note: CHITEST [Compatibility] Returns the test for independence
	// Note: CHISQ.DIST [Statistical] Returns the cumulative beta probability density
	// CHISQ.DIST.RT [Statistical] Returns the one-tailed probability of the chi-squared distribution
	// CHISQ.INV [Statistical] Returns the cumulative beta probability density
	// CHISQ.INV.RT [Statistical] Returns the inverse of the one-tailed probability of the chi-squared distribution
	// CHISQ.TEST [Statistical] Returns the test for independence
	// CHOOSE [Lookup and reference] Chooses a value from a list of values
	// CLEAN [Text] Removes all nonprintable characters from text
	// CODE [Text] Returns a numeric code for the first character in a text string
	// COLUMN [Lookup and reference] Returns the column number of a reference
	// COLUMNS [Lookup and reference] Returns the number of columns in a reference
	// COMBIN [Math and trigonometry] Returns the number of combinations for a given number of objects
	// COMBINA [Math and trigonometry:   ] Returns the number of combinations with repetitions for a given number of items
	// COMPLEX [Engineering] Converts real and imaginary coefficients into a complex number
	// CONCAT [Text] Combines the text from multiple ranges and/or strings, but it doesn't provide the delimiter or IgnoreEmpty arguments.
	// CONCATENATE [Text] Joins several text items into one text item
	// CONFIDENCE [Compatibility] Returns the confidence interval for a population mean
	// CONFIDENCE.NORM [Statistical] Returns the confidence interval for a population mean
	// CONFIDENCE.T [Statistical] Returns the confidence interval for a population mean, using a Student's t distribution
	// CONVERT [Engineering] Converts a number from one measurement system to another
	// CORREL [Statistical] Returns the correlation coefficient between two data sets
	// COS [Math and trigonometry] Returns the cosine of a number
	// COSH [Math and trigonometry] Returns the hyperbolic cosine of a number
	// COT [Math and trigonometry] Returns the hyperbolic cosine of a number
	// COTH [Math and trigonometry] Returns the cotangent of an angle
	// COUNT [Statistical] Counts how many numbers are in the list of arguments
	// COUNTA [Statistical] Counts how many values are in the list of arguments
	// COUNTBLANK [Statistical] Counts the number of blank cells within a range
	// COUNTIF [Statistical] Counts the number of cells within a range that meet the given criteria
	// COUNTIFS [Statistical] Counts the number of cells within a range that meet multiple criteria
	// COUPDAYBS [Financial] Returns the number of days from the beginning of the coupon period to the settlement date
	// COUPDAYS [Financial] Returns the number of days in the coupon period that contains the settlement date
	// COUPDAYSNC [Financial] Returns the number of days from the settlement date to the next coupon date
	// COUPNCD [Financial] Returns the next coupon date after the settlement date
	// COUPNUM [Financial] Returns the number of coupons payable between the settlement date and maturity date
	// COUPPCD [Financial] Returns the previous coupon date before the settlement date
	// COVAR [Compatibility] Returns covariance, the average of the products of paired deviations
	// COVARIANCE.P [Statistical] Returns covariance, the average of the products of paired deviations
	// COVARIANCE.S [Statistical] Returns the sample covariance, the average of the products deviations for each data point pair in two data sets
	// CRITBINOM [Compatibility] Returns the smallest value for which the cumulative binomial distribution is less than or equal to a criterion value
	// CSC [Math and trigonometry] Returns the cosecant of an angle
	// CSCH [Math and trigonometry] Returns the hyperbolic cosecant of an angle
	// CUBEKPIMEMBER [Cube] Returns a key performance indicator (KPI) name, property, and measure, and displays the name and property in the cell. A KPI is a quantifiable measurement, such as monthly gross profit or quarterly employee turnover, used to monitor an organization's performance.
	// CUBEMEMBER [Cube] Returns a member or tuple in a cube hierarchy. Use to validate that the member or tuple exists in the cube.
	// CUBEMEMBERPROPERTY [Cube] Returns the value of a member property in the cube. Use to validate that a member name exists within the cube and to return the specified property for this member.
	// CUBERANKEDMEMBER [Cube] Returns the nth, or ranked, member in a set. Use to return one or more elements in a set, such as the top sales performer or top 10 students.
	// CUBESET [Cube] Defines a calculated set of members or tuples by sending a set expression to the cube on the server, which creates the set, and then returns that set to Microsoft Office Excel.
	// CUBESETCOUNT [Cube] Returns the number of items in a set.
	// CUBEVALUE [Cube] Returns an aggregated value from a cube.
	// CUMIPMT [Financial] Returns the cumulative interest paid between two periods
	// CUMPRINC [Financial] Returns the cumulative principal paid on a loan between two periods
	// DATE [Date and time] Returns the serial number of a particular date
	// DATEDIF [Date and time] Calculates the number of days, months, or years between two dates. This function is useful in formulas where you need to calculate an age.
	// DATEVALUE [Date and time] Converts a date in the form of text to a serial number
	// DAVERAGE [Database] Returns the average of selected database entries
	// DAY [Date and time] Converts a serial number to a day of the month
	// DAYS [Date and time] Returns the number of days between two dates
	// DAYS360 [Date and time] Calculates the number of days between two dates based on a 360-day year
	// DB [Financial] Returns the depreciation of an asset for a specified period by using the fixed-declining balance method
	// DBCS [Text] Changes half-width (single-byte) English letters or katakana within a character string to full-width (double-byte) characters
	// DCOUNT [Database] Counts the cells that contain numbers in a database
	// DCOUNTA [Database] Counts nonblank cells in a database
	// DDB [Financial] Returns the depreciation of an asset for a specified period by using the double-declining balance method or some other method that you specify
	// DEC2BIN [Engineering] Converts a decimal number to binary
	// DEC2HEX [Engineering] Converts a decimal number to hexadecimal
	// DEC2OCT [Engineering] Converts a decimal number to octal
	// DECIMAL [Math and trigonometry] Converts a text representation of a number in a given base into a decimal number
	// DEGREES [Math and trigonometry] Converts radians to degrees
	// DELTA [Engineering] Tests whether two values are equal
	// DEVSQ [Statistical] Returns the sum of squares of deviations
	// DGET [Database] Extracts from a database a single record that matches the specified criteria
	// DISC [Financial] Returns the discount rate for a security
	// DMAX [Database] Returns the maximum value from selected database entries
	// DMIN [Database] Returns the minimum value from selected database entries
	// DOLLAR [Text] Converts a number to text, using the $ (dollar) currency format
	// DOLLARDE [Financial] Converts a dollar price, expressed as a fraction, into a dollar price, expressed as a decimal number
	// DOLLARFR [Financial] Converts a dollar price, expressed as a decimal number, into a dollar price, expressed as a fraction
	// DPRODUCT [Database] Multiplies the values in a particular field of records that match the criteria in a database
	// DSTDEV [Database] Estimates the standard deviation based on a sample of selected database entries
	// DSTDEVP [Database] Calculates the standard deviation based on the entire population of selected database entries
	// DSUM [Database] Adds the numbers in the field column of records in the database that match the criteria
	// DURATION [Financial] Returns the annual duration of a security with periodic interest payments
	// DVAR [Database] Estimates variance based on a sample from selected database entries
	// DVARP [Database] Calculates variance based on the entire population of selected database entries
	// EDATE [Date and time] Returns the serial number of the date that is the indicated number of months before or after the start date
	// EFFECT [Financial] Returns the effective annual interest rate
	// ENCODEURL [Web] Returns a URL-encoded string
	// EOMONTH [Date and time] Returns the serial number of the last day of the month before or after a specified number of months
	// ERF [Engineering] Returns the error
	// ERF.PRECISE [Engineering] Returns the error
	// ERFC [Engineering] Returns the complementary error
	// ERFC.PRECISE [Engineering] Returns the complementary ERF function integrated between x and infinity
	// ERROR.TYPE [Information] Returns a number corresponding to an error type
	// EUROCONVERT [Add-in and Automation] Converts a number to euros, converts a number from euros to a euro member currency, or converts a number from one euro member currency to another by using the euro as an intermediary (triangulation).
	// EVEN [Math and trigonometry] Rounds a number up to the nearest even integer
	// EXACT [Text] Checks to see if two text values are identical
	// EXP [Math and trigonometry] Returns e raised to the power of a given number
	// EXPON.DIST [Statistical] Returns the exponential distribution
	// EXPONDIST [Compatibility] Returns the exponential distribution
	// FACT [Math and trigonometry] Returns the factorial of a number
	// FACTDOUBLE [Math and trigonometry] Returns the double factorial of a number
	// FALSE [Logical] Returns the logical value FALSE
	// F.DIST [Statistical] Returns the F probability distribution
	// FDIST [Compatibility] Returns the F probability distribution
	// F.DIST.RT [Statistical] Returns the F probability distribution
	// FILTER [Lookup and reference] Filters a range of data based on criteria you define
	// FILTERXML [Web] Returns specific data from the XML content by using the specified XPath
	// FIND, FINDB [Text] Finds one text value within another (case-sensitive)
	// F.INV [Statistical] Returns the inverse of the F probability distribution
	// F.INV.RT [Statistical] Returns the inverse of the F probability distribution
	// FINV [Statistical] Returns the inverse of the F probability distribution
	// FISHER [Statistical] Returns the Fisher transformation
	// FISHERINV [Statistical] Returns the inverse of the Fisher transformation
	// FIXED [Text] Formats a number as text with a fixed number of decimals
	// FLOOR [Compatibility] Rounds a number down, toward zero
	// FLOOR.MATH [Math and trigonometry] Rounds a number down, to the nearest integer or to the nearest multiple of significance
	// FLOOR.PRECISE [Math and trigonometry] Rounds a number the nearest integer or to the nearest multiple of significance. Regardless of the sign of the number, the number is rounded up.
	// FORECAST [Statistical] Returns a value along a linear trend
	// FORECAST.ETS [Statistical] Returns a future value based on existing (historical) values by using the AAA version of the Exponential Smoothing (ETS) algorithm
	// FORECAST.ETS.CONFINT [Statistical] Returns a confidence interval for the forecast value at the specified target date
	// FORECAST.ETS.SEASONALITY [Statistical] Returns the length of the repetitive pattern Excel detects for the specified time series
	// FORECAST.ETS.STAT [Statistical] Returns a statistical value as a result of time series forecasting
	// FORECAST.LINEAR [Statistical] Returns a future value based on existing values
	// FORMULATEXT [Lookup and reference] Returns the formula at the given reference as text
	// FREQUENCY [Statistical] Returns a frequency distribution as a vertical array
	// F.TEST [Statistical] Returns the result of an F-test
	// FTEST [Compatibility] Returns the result of an F-test
	// FV [Financial] Returns the future value of an investment
	// FVSCHEDULE [Financial] Returns the future value of an initial principal after applying a series of compound interest rates
	// GAMMA [Statistical] Returns the Gamma function value
	// GAMMA.DIST [Statistical] Returns the gamma distribution
	// GAMMADIST [Compatibility] Returns the gamma distribution
	// GAMMA.INV [Statistical] Returns the inverse of the gamma cumulative distribution
	// GAMMAINV [Compatibility] Returns the inverse of the gamma cumulative distribution
	// GAMMALN [Statistical] Returns the natural logarithm of the gamma function, Γ(x)
	// GAMMALN.PRECISE [Statistical] Returns the natural logarithm of the gamma function, Γ(x)
	// GAUSS [Statistical] Returns 0.5 less than the standard normal cumulative distribution
	// GCD [Math and trigonometry] Returns the greatest common divisor
	// GEOMEAN [Statistical] Returns the geometric mean
	// GESTEP [Engineering] Tests whether a number is greater than a threshold value
	// GETPIVOTDATA [Lookup and reference] Returns data stored in a PivotTable report
	// GROWTH [Statistical] Returns values along an exponential trend
	// HARMEAN [Statistical] Returns the harmonic mean
	// HEX2BIN [Engineering] Converts a hexadecimal number to binary
	// HEX2DEC [Engineering] Converts a hexadecimal number to decimal
	// HEX2OCT [Engineering] Converts a hexadecimal number to octal
	// HLOOKUP [Lookup and reference] Looks in the top row of an array and returns the value of the indicated cell
	// HOUR [Date and time] Converts a serial number to an hour
	// HYPERLINK [Lookup and reference] Creates a shortcut or jump that opens a document stored on a network server, an intranet, or the Internet
	// HYPGEOM.DIST [Statistical] Returns the hypergeometric distribution
	// HYPGEOMDIST [Compatibility] Returns the hypergeometric distribution
	// IFERROR [Logical] Returns a value you specify if a formula evaluates to an error; otherwise, returns the result of the formula
	// IFNA [Logical] Returns the value you specify if the expression resolves to #N/A, otherwise returns the result of the expression
	// IFS [Logical] Checks whether one or more conditions are met and returns a value that corresponds to the first TRUE condition.
	// IMABS [Engineering] Returns the absolute value (modulus) of a complex number
	// IMAGINARY [Engineering] Returns the imaginary coefficient of a complex number
	// IMARGUMENT [Engineering] Returns the argument theta, an angle expressed in radians
	// IMCONJUGATE [Engineering] Returns the complex conjugate of a complex number
	// IMCOS [Engineering] Returns the cosine of a complex number
	// IMCOSH [Engineering] Returns the hyperbolic cosine of a complex number
	// IMCOT [Engineering] Returns the cotangent of a complex number
	// IMCSC [Engineering] Returns the cosecant of a complex number
	// IMCSCH [Engineering] Returns the hyperbolic cosecant of a complex number
	// IMDIV [Engineering] Returns the quotient of two complex numbers
	// IMEXP [Engineering] Returns the exponential of a complex number
	// IMLN [Engineering] Returns the natural logarithm of a complex number
	// IMLOG10 [Engineering] Returns the base-10 logarithm of a complex number
	// IMLOG2 [Engineering] Returns the base-2 logarithm of a complex number
	// IMPOWER [Engineering] Returns a complex number raised to an integer power
	// IMPRODUCT [Engineering] Returns the product of complex numbers
	// IMREAL [Engineering] Returns the real coefficient of a complex number
	// IMSEC [Engineering] Returns the secant of a complex number
	// IMSECH [Engineering] Returns the hyperbolic secant of a complex number
	// IMSIN [Engineering] Returns the sine of a complex number
	// IMSINH [Engineering] Returns the hyperbolic sine of a complex number
	// IMSQRT [Engineering] Returns the square root of a complex number
	// IMSUB [Engineering] Returns the difference between two complex numbers
	// IMSUM [Engineering] Returns the sum of complex numbers
	// IMTAN [Engineering] Returns the tangent of a complex number
	// INDEX [Lookup and reference] Uses an index to choose a value from a reference or array
	// INDIRECT [Lookup and reference] Returns a reference indicated by a text value
	// INFO [Information] Returns information about the current operating environment
	// INT [Math and trigonometry] Rounds a number down to the nearest integer
	// INTERCEPT [Statistical] Returns the intercept of the linear regression line
	// INTRATE [Financial] Returns the interest rate for a fully invested security
	// IPMT [Financial] Returns the interest payment for an investment for a given period
	// IRR [Financial] Returns the internal rate of return for a series of cash flows
	// ISBLANK [Information] Returns TRUE if the value is blank
	// ISERR [Information] Returns TRUE if the value is any error value except #N/A
	// ISERROR [Information] Returns TRUE if the value is any error value
	// ISEVEN [Information] Returns TRUE if the number is even
	// ISFORMULA [Information] Returns TRUE if there is a reference to a cell that contains a formula
	// ISLOGICAL [Information] Returns TRUE if the value is a logical value
	// ISNA [Information] Returns TRUE if the value is the #N/A error value
	// ISNONTEXT [Information] Returns TRUE if the value is not text
	// ISNUMBER [Information] Returns TRUE if the value is a number
	// ISODD [Information] Returns TRUE if the number is odd
	// ISREF [Information] Returns TRUE if the value is a reference
	// ISTEXT [Information] Returns TRUE if the value is text
	// ISO.CEILING [Math and trigonometry] Returns a number that is rounded up to the nearest integer or to the nearest multiple of significance
	// ISOWEEKNUM [Date and time] Returns the number of the ISO week number of the year for a given date
	// ISPMT [Financial] Calculates the interest paid during a specific period of an investment
	// JIS [Text] Changes half-width (single-byte) characters within a string to full-width (double-byte) characters
	// KURT [Statistical] Returns the kurtosis of a data set
	// LARGE [Statistical] Returns the k-th largest value in a data set
	// LCM [Math and trigonometry] Returns the least common multiple
	// LEFT, LEFTB [Text] Returns the leftmost characters from a text value
	// LEN, LENB [Text] Returns the number of characters in a text string
	// LINEST [Statistical] Returns the parameters of a linear trend
	// LN [Math and trigonometry] Returns the natural logarithm of a number
	// LOG [Math and trigonometry] Returns the logarithm of a number to a specified base
	// LOG10 [Math and trigonometry] Returns the base-10 logarithm of a number
	// LOGEST [Statistical] Returns the parameters of an exponential trend
	// LOGINV [Compatibility] Returns the inverse of the lognormal cumulative distribution
	// LOGNORM.DIST [Statistical] Returns the cumulative lognormal distribution
	// LOGNORMDIST [Compatibility] Returns the cumulative lognormal distribution
	// LOGNORM.INV [Statistical] Returns the inverse of the lognormal cumulative distribution
	// LOOKUP [Lookup and reference] Looks up values in a vector or array
	// LOWER [Text] Converts text to lowercase
	// MATCH [Lookup and reference] Looks up values in a reference or array
	// MAX [Statistical] Returns the maximum value in a list of arguments
	// MAXA [Statistical] Returns the maximum value in a list of arguments, including numbers, text, and logical values
	// MAXIFS [Statistical] Returns the maximum value among cells specified by a given set of conditions or criteria
	// MDETERM [Math and trigonometry] Returns the matrix determinant of an array
	// MDURATION [Financial] Returns the Macauley modified duration for a security with an assumed par value of $100
	// MEDIAN [Statistical] Returns the median of the given numbers
	// MID, MIDB [Text] Returns a specific number of characters from a text string starting at the position you specify
	// MIN [Statistical] Returns the minimum value in a list of arguments
	// MINIFS [Statistical] Returns the minimum value among cells specified by a given set of conditions or criteria.
	// MINA [Statistical] Returns the smallest value in a list of arguments, including numbers, text, and logical values
	// MINUTE [Date and time] Converts a serial number to a minute
	// MINVERSE [Math and trigonometry] Returns the matrix inverse of an array
	// MIRR [Financial] Returns the internal rate of return where positive and negative cash flows are financed at different rates
	// MMULT [Math and trigonometry] Returns the matrix product of two arrays
	// MOD [Math and trigonometry] Returns the remainder from division
	// MODE [Compatibility] Returns the most common value in a data set
	// MODE.MULT [Statistical] Returns a vertical array of the most frequently occurring, or repetitive values in an array or range of data
	// MODE.SNGL [Statistical] Returns the most common value in a data set
	// MONTH [Date and time] Converts a serial number to a month
	// MROUND [Math and trigonometry] Returns a number rounded to the desired multiple
	// MULTINOMIAL [Math and trigonometry] Returns the multinomial of a set of numbers
	// MUNIT [Math and trigonometry] Returns the unit matrix or the specified dimension
	// N [Information] Returns a value converted to a number
	// NA [Information] Returns the error value #N/A
	// NEGBINOM.DIST [Statistical] Returns the negative binomial distribution
	// NEGBINOMDIST [Compatibility] Returns the negative binomial distribution
	// NETWORKDAYS [Date and time] Returns the number of whole workdays between two dates
	// NETWORKDAYS.INTL [Date and time] Returns the number of whole workdays between two dates using parameters to indicate which and how many days are weekend days
	// NOMINAL [Financial] Returns the annual nominal interest rate
	// NORM.DIST [Statistical] Returns the normal cumulative distribution
	// NORMDIST [Compatibility] Returns the normal cumulative distribution
	// NORMINV [Statistical] Returns the inverse of the normal cumulative distribution
	// NORM.INV [Compatibility] Returns the inverse of the normal cumulative distribution
	// Note: NORM.S.DIST [Statistical] Returns the standard normal cumulative distribution
	// NORMSDIST [Compatibility] Returns the standard normal cumulative distribution
	// NORM.S.INV [Statistical] Returns the inverse of the standard normal cumulative distribution
	// NORMSINV [Compatibility] Returns the inverse of the standard normal cumulative distribution
	// NOT [Logical] Reverses the logic of its argument
	// NOW [Date and time] Returns the serial number of the current date and time
	// NPER [Financial] Returns the number of periods for an investment
	// NPV [Financial] Returns the net present value of an investment based on a series of periodic cash flows and a discount rate
	// NUMBERVALUE [Text] Converts text to number in a locale-independent manner
	// OCT2BIN [Engineering] Converts an octal number to binary
	// OCT2DEC [Engineering] Converts an octal number to decimal
	// OCT2HEX [Engineering] Converts an octal number to hexadecimal
	// ODD [Math and trigonometry] Rounds a number up to the nearest odd integer
	// ODDFPRICE [Financial] Returns the price per $100 face value of a security with an odd first period
	// ODDFYIELD [Financial] Returns the yield of a security with an odd first period
	// ODDLPRICE [Financial] Returns the price per $100 face value of a security with an odd last period
	// ODDLYIELD [Financial] Returns the yield of a security with an odd last period
	// OFFSET [Lookup and reference] Returns a reference offset from a given reference
	// OR [Logical] Returns TRUE if any argument is TRUE
	// PDURATION [Financial] Returns the number of periods required by an investment to reach a specified value
	// PEARSON [Statistical] Returns the Pearson product moment correlation coefficient
	// PERCENTILE.EXC [Statistical] Returns the k-th percentile of values in a range, where k is in the range 0..1, exclusive
	// PERCENTILE.INC [Statistical] Returns the k-th percentile of values in a range
	// PERCENTILE [Compatibility] Returns the k-th percentile of values in a range
	// PERCENTRANK.EXC [Statistical] Returns the rank of a value in a data set as a percentage (0..1, exclusive) of the data set
	// PERCENTRANK.INC [Statistical] Returns the percentage rank of a value in a data set
	// PERCENTRANK [Compatibility] Returns the percentage rank of a value in a data set
	// PERMUT [Statistical] Returns the number of permutations for a given number of objects
	// PERMUTATIONA [Statistical] Returns the number of permutations for a given number of objects (with repetitions) that can be selected from the total objects
	// PHI [Statistical] Returns the value of the density function for a standard normal distribution
	// PHONETIC [Text] Extracts the phonetic (furigana) characters from a text string
	// PI [Math and trigonometry] Returns the value of pi
	// PMT [Financial] Returns the periodic payment for an annuity
	// POISSON.DIST [Statistical] Returns the Poisson distribution
	// POISSON [Compatibility] Returns the Poisson distribution
	// POWER [Math and trigonometry] Returns the result of a number raised to a power
	// PPMT [Financial] Returns the payment on the principal for an investment for a given period
	// PRICE [Financial] Returns the price per $100 face value of a security that pays periodic interest
	// PRICEDISC [Financial] Returns the price per $100 face value of a discounted security
	// PRICEMAT [Financial] Returns the price per $100 face value of a security that pays interest at maturity
	// PROB [Statistical] Returns the probability that values in a range are between two limits
	// PRODUCT [Math and trigonometry] Multiplies its arguments
	// PROPER [Text] Capitalizes the first letter in each word of a text value
	// PV [Financial] Returns the present value of an investment
	// QUARTILE [Compatibility] Returns the quartile of a data set
	// QUARTILE.EXC [Statistical] Returns the quartile of the data set, based on percentile values from 0..1, exclusive
	// QUARTILE.INC [Statistical] Returns the quartile of a data set
	// QUOTIENT [Math and trigonometry] Returns the integer portion of a division
	// RADIANS [Math and trigonometry] Converts degrees to radians
	// RAND [Math and trigonometry] Returns a random number between 0 and 1
	// RANDARRAY [Math and trigonometry] Returns an array of random numbers between 0 and 1
	// RANDBETWEEN [Math and trigonometry] Returns a random number between the numbers you specify
	// RANK.AVG [Statistical] Returns the rank of a number in a list of numbers
	// RANK.EQ [Statistical] Returns the rank of a number in a list of numbers
	// RANK [Compatibility] Returns the rank of a number in a list of numbers
	// RATE [Financial] Returns the interest rate per period of an annuity
	// RECEIVED [Financial] Returns the amount received at maturity for a fully invested security
	// REGISTER.ID [Add-in and Automation] Returns the register ID of the specified dynamic link library (DLL) or code resource that has been previously registered
	// REPLACE, REPLACEB [Text] Replaces characters within text
	// REPT [Text] Repeats text a given number of times
	// RIGHT, RIGHTB [Text] Returns the rightmost characters from a text value
	// ROMAN [Math and trigonometry] Converts an arabic numeral to roman, as text
	// ROUND [Math and trigonometry] Rounds a number to a specified number of digits
	// ROUNDDOWN [Math and trigonometry] Rounds a number down, toward zero
	// ROUNDUP [Math and trigonometry] Rounds a number up, away from zero
	// ROW [Lookup and reference] Returns the row number of a reference
	// ROWS [Lookup and reference] Returns the number of rows in a reference
	// RRI [Financial] Returns an equivalent interest rate for the growth of an investment
	// RSQ [Statistical] Returns the square of the Pearson product moment correlation coefficient
	// RTD [Lookup and reference] Retrieves real-time data from a program that supports COM automation
	// SEARCH, SEARCHB [Text] Finds one text value within another (not case-sensitive)
	// SEC [Math and trigonometry] Returns the secant of an angle
	// SECH [Math and trigonometry] Returns the hyperbolic secant of an angle
	// SECOND [Date and time] Converts a serial number to a second
	// SEQUENCE [Math and trigonometry] Generates a list of sequential numbers in an array, such as 1, 2, 3, 4
	// SERIESSUM [Math and trigonometry] Returns the sum of a power series based on the formula
	// SHEET [Information] Returns the sheet number of the referenced sheet
	// SHEETS [Information] Returns the number of sheets in a reference
	// SIGN [Math and trigonometry] Returns the sign of a number
	// SIN [Math and trigonometry] Returns the sine of the given angle
	// SINGLE [Lookup and reference] Returns a single value using logic known as implicit intersection
	// SINH [Math and trigonometry] Returns the hyperbolic sine of a number
	// SKEW [Statistical] Returns the skewness of a distribution
	// SKEW.P [Statistical] Returns the skewness of a distribution based on a population: a characterization of the degree of asymmetry of a distribution around its mean
	// SLN [Financial] Returns the straight-line depreciation of an asset for one period
	// SLOPE [Statistical] Returns the slope of the linear regression line
	// SMALL [Statistical] Returns the k-th smallest value in a data set
	// SORT [Lookup and reference] Sorts the contents of a range or array
	// SORTBY [Lookup and reference] Sorts the contents of a range or array based on the values in a corresponding range or array
	// SQRT [Math and trigonometry] Returns a positive square root
	// SQRTPI [Math and trigonometry] Returns the square root of (number * pi)
	// STANDARDIZE [Statistical] Returns a normalized value
	// STDEV [Compatibility] Estimates standard deviation based on a sample
	// STDEV.P [Statistical] Calculates standard deviation based on the entire population
	// STDEV.S [Statistical] Estimates standard deviation based on a sample
	// STDEVA [Statistical] Estimates standard deviation based on a sample, including numbers, text, and logical values
	// STDEVP [Compatibility] Calculates standard deviation based on the entire population
	// STDEVPA [Statistical] Calculates standard deviation based on the entire population, including numbers, text, and logical values
	// STEYX [Statistical] Returns the standard error of the predicted y-value for each x in the regression
	// SUBSTITUTE [Text] Substitutes new text for old text in a text string
	// SUBTOTAL [Math and trigonometry] Returns a subtotal in a list or database
	// SUM [Math and trigonometry] Adds its arguments
	// SUMIF [Math and trigonometry] Adds the cells specified by a given criteria
	// SUMIFS [Math and trigonometry] Adds the cells in a range that meet multiple criteria
	// SUMPRODUCT [Math and trigonometry] Returns the sum of the products of corresponding array components
	// SUMSQ [Math and trigonometry] Returns the sum of the squares of the arguments
	// SUMX2MY2 [Math and trigonometry] Returns the sum of the difference of squares of corresponding values in two arrays
	// SUMX2PY2 [Math and trigonometry] Returns the sum of the sum of squares of corresponding values in two arrays
	// SUMXMY2 [Math and trigonometry] Returns the sum of squares of differences of corresponding values in two arrays
	// SWITCH [Logical] Evaluates an expression against a list of values and returns the result corresponding to the first matching value. If there is no match, an optional default value may be returned.
	// SYD [Financial] Returns the sum-of-years' digits depreciation of an asset for a specified period
	// T [Text] Converts its arguments to text
	// TAN [Math and trigonometry] Returns the tangent of a number
	// TANH [Math and trigonometry] Returns the hyperbolic tangent of a number
	// TBILLEQ [Financial] Returns the bond-equivalent yield for a Treasury bill
	// TBILLPRICE [Financial] Returns the price per $100 face value for a Treasury bill
	// TBILLYIELD [Financial] Returns the yield for a Treasury bill
	// T.DIST [Statistical] Returns the Percentage Points (probability) for the Student t-distribution
	// T.DIST.2T [Statistical] Returns the Percentage Points (probability) for the Student t-distribution
	// T.DIST.RT [Statistical] Returns the Student's t-distribution
	// TDIST [Compatibility] Returns the Student's t-distribution
	// TEXT [Text] Formats a number and converts it to text
	// TEXTJOIN [Text] Combines the text from multiple ranges and/or strings, and includes a delimiter you specify between each text value that will be combined. If the delimiter is an empty text string, this function will effectively concatenate the ranges.
	// TIME [Date and time] Returns the serial number of a particular time
	// TIMEVALUE [Date and time] Converts a time in the form of text to a serial number
	// T.INV [Statistical] Returns the t-value of the Student's t-distribution as a function of the probability and the degrees of freedom
	// T.INV.2T [Statistical] Returns the inverse of the Student's t-distribution
	// TINV [Compatibility] Returns the inverse of the Student's t-distribution
	// TODAY [Date and time] Returns the serial number of today's date
	// TRANSPOSE [Lookup and reference] Returns the transpose of an array
	// TREND [Statistical] Returns values along a linear trend
	// TRIM [Text] Removes spaces from text
	// TRIMMEAN [Statistical] Returns the mean of the interior of a data set
	// TRUE [Logical] Returns the logical value TRUE
	// TRUNC [Math and trigonometry] Truncates a number to an integer
	// T.TEST [Statistical] Returns the probability associated with a Student's t-test
	// TTEST [Compatibility] Returns the probability associated with a Student's t-test
	// TYPE [Information] Returns a number indicating the data type of a value
	// UNICHAR [Text] Returns the Unicode character that is references by the given numeric value
	// UNICODE [Text] Returns the number (code point) that corresponds to the first character of the text
	// UNIQUE [Lookup and reference] Returns a list of unique values in a list or range
	// UPPER [Text] Converts text to uppercase
	// VALUE [Text] Converts a text argument to a number
	// VAR [Compatibility] Estimates variance based on a sample
	// VAR.P [Statistical] Calculates variance based on the entire population
	// VAR.S [Statistical] Estimates variance based on a sample
	// VARA [Statistical] Estimates variance based on a sample, including numbers, text, and logical values
	// VARP [Compatibility] Calculates variance based on the entire population
	// VARPA [Statistical] Calculates variance based on the entire population, including numbers, text, and logical values
	// VDB [Financial] Returns the depreciation of an asset for a specified or partial period by using a declining balance method
	// VLOOKUP [Lookup and reference] Looks in the first column of an array and moves across the row to return the value of a cell
	// WEBSERVICE [Web] Returns data from a web service.
	// WEEKDAY [Date and time] Converts a serial number to a day of the week
	// WEEKNUM [Date and time] Converts a serial number to a number representing where the week falls numerically with a year
	// WEIBULL [Compatibility] Calculates variance based on the entire population, including numbers, text, and logical values
	// WEIBULL.DIST [Statistical] Returns the Weibull distribution
	// WORKDAY [Date and time] Returns the serial number of the date before or after a specified number of workdays
	// WORKDAY.INTL [Date and time] Returns the serial number of the date before or after a specified number of workdays using parameters to indicate which and how many days are weekend days
	// XIRR [Financial] Returns the internal rate of return for a schedule of cash flows that is not necessarily periodic
	// XNPV [Financial] Returns the net present value for a schedule of cash flows that is not necessarily periodic
	// XOR [Logical] Returns a logical exclusive OR of all arguments
	// YEAR [Date and time] Converts a serial number to a year
	// YEARFRAC [Date and time] Returns the year fraction representing the number of whole days between start_date and end_date
	// YIELD [Financial] Returns the yield on a security that pays periodic interest
	// YIELDDISC [Financial] Returns the annual yield for a discounted security; for example, a Treasury bill
	// YIELDMAT [Financial] Returns the annual yield of a security that pays interest at maturity
	// Z.TEST [Statistical] Returns the one-tailed probability-value of a z-test
	// ZTEST [Compatibility] Returns the one-tailed probability-value of a z-test
}

func trim(ec *eval.Context, args []eval.Value) (eval.Value, error) {
	s, err := args[0].StringValue(ec)
	if err != nil {
		return eval.NewEmptyValue(), err
	}
	return eval.NewStringValue(strings.Trim(s, "\n\r\t ")), nil
}

func sum(ec *eval.Context, args []eval.Value) (eval.Value, error) {
	s := decimal.Zero
	for i := range args {
		if args[i].Type() == eval.TypeRangeRef {
			err := ec.DataProvider.IterateDecimalValues(
				ec,
				args[i].Cell().CellAddress,
				args[i].CellTo().CellAddress,
				func(d decimal.Decimal) error {
					s = s.Add(d)
					return nil
				},
			)
			if err != nil {
				return eval.NewEmptyValue(), err
			}
		} else {
			d, err := args[i].DecimalValue(ec)
			if err != nil {
				return eval.NewEmptyValue(), err
			}
			s = s.Add(d)
		}
	}
	return eval.NewDecimalValue(s), nil
}

// IF [Logical] Specifies a logical test to perform
func if_(ec *eval.Context, args []eval.Value) (eval.Value, error) {
	b, err := args[0].BoolValue(ec)
	if err != nil {
		return eval.NewEmptyValue(), err
	}
	if b {
		return args[1], nil
	} else {
		return args[2], nil
	}
}
