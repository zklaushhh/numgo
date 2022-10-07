package numgo

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"

	"gocv.io/x/gocv"
)

type Matrix struct {
	M   []float32
	Row uint
	Col uint
}

// 计算矩阵平均值
func (m Matrix) Mean() float32 {
	var sum float32
	for i := 0; i < int(m.Row*m.Col); i++ {
		sum = sum + m.M[i]
	}
	result := sum / float32(m.Row*m.Col)
	return result
}

// 复制矩阵
func Copy(m Matrix) Matrix {
	var matrix Matrix
	matrix.Row = m.Row
	matrix.Col = m.Col

	M := make([]float32, m.Row*m.Col)
	for i := 0; i < int(m.Row*m.Col); i++ {
		M[i] = m.M[i]
	}
	matrix.M = M
	return matrix
}

// 将矩阵绝对值化
func Abs(m Matrix) Matrix {
	var matrix Matrix
	matrix.Row = m.Row
	matrix.Col = m.Col

	M := make([]float32, m.Row*m.Col)
	for i := 0; i < int(m.Row*m.Col); i++ {
		if m.M[i] >= 0 {
			M[i] = m.M[i]
		} else {
			M[i] = m.M[i] * (-1)
		}
	}

	matrix.M = M

	return matrix

}

// 生成全是1的对角矩阵(行列相等)
func Diag(size uint) Matrix {

	matrix := Zeros(size, size)

	var p int
	for i := 0; i < int(size); i++ {
		matrix.M[i*int(size)+p] = 1
		p++
	}

	return matrix

}

// 将矩阵转换为一维的[]byte
func ConvertToBytes(m Matrix) []byte {
	var result []byte
	for i := 0; i < int(m.Row*m.Col); i++ {
		result = append(result, byte(m.M[i]))
	}
	return result

}

// 返回数组中最大的值
func (m Matrix) Max() float32 {
	max := m.M[0]
	for i := 1; i < int(m.Row*m.Col); i++ {
		if m.M[i] > max {
			max = m.M[i]
		}
	}
	return max

}

// 返回切片里非空下标
func SliceNotNull(s []float32) []int {
	if len(s) == 0 {
		return []int{}
	}

	var result []int
	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			result = append(result, i)
		}
	}
	return result

}

// 返回类似np.where(v_filter>thresh_e)的效果
// 比较对象为两个矩阵
func MatrixWhere(m1, m2 Matrix, flag string) ([]int, error) {
	// 只适用于列数为1的矩阵
	if m1.Col != 1 || m2.Col != 1 {
		return []int{}, fmt.Errorf("invalid matrix parameter")
	}

	if !(flag == ">" || flag == "<" || flag == ">=" || flag == "<=") {
		return []int{}, fmt.Errorf("invalid flag parameter")
	}

	var result []int

	switch flag {
	case ">":
		// 返回满足(m > value)的下标
		for i := 0; i < int(m1.Row); i++ {
			if m1.M[i] > m2.M[i] {
				result = append(result, i)
			}
		}
	case ">=":
		// 返回满足(m >= value)的下标
		for i := 0; i < int(m1.Row); i++ {
			if m1.M[i] >= m2.M[i] {
				result = append(result, i)
			}
		}

	case "<":
		// 返回满足(m < value)的下标
		for i := 0; i < int(m1.Row); i++ {
			if m1.M[i] < m2.M[i] {
				result = append(result, i)
			}
		}

	case "<=":
		// 返回满足(m <= value)的下标
		for i := 0; i < int(m1.Row); i++ {
			if m1.M[i] <= m2.M[i] {
				result = append(result, i)
			}
		}
	}

	return result, nil

}

// 将n行1列的数组 按列填充 到n行m列的数组中
// m2为要插入的元素, j指定插入到m的第几列
func (m *Matrix) Insert(m2 Matrix, j int) error {
	if m.Row != m2.Row || m2.Col != 1 {
		return fmt.Errorf("invalid parameter")
	}

	for i := 0; i < int(m.Row); i++ {
		m.M[i*int(m.Col)+j] = m2.M[i]
	}

	return nil

}

// 用切片s的内容填充数组某列
func (m *Matrix) InsertBySlice(s []float32, j int) error {
	if len(s) != int(m.Row) {
		return fmt.Errorf("INVALID PARAMETER")
	}
	for i := 0; i < len(s); i++ {
		m.M[i*int(m.Col)+j] = s[i]
	}
	return nil
}

// 矩阵对应位置相乘
func MatrixMultiply(m1, m2 Matrix) Matrix {
	if (m1.Row != m2.Row) || (m1.Col != m2.Col) {
		return Matrix{}
	}

	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)
	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] * m2.M[i]

	}
	matrix.M = M
	return matrix
}

// 数值减矩阵
// 先把数值转换成矩阵, 和矩阵相减后返回一个矩阵
func NumSubMatrix(m Matrix, value float32) Matrix {
	matrix := Ones(m.Row, m.Col)
	matrix.Multiply(value)

	for i := 0; i < int(m.Row*m.Col); i++ {
		matrix.M[i] = matrix.M[i] - m.M[i]
	}

	return matrix

}

// 根据所给的slice，取不为0的部分的值填充到数组对应位置中
func (m *Matrix) SetValueBySlice(m2 Matrix, s []float32) error {
	if (int(m2.Row) != len(s)) || m2.Col != 1 || m.Row != m2.Row {
		return fmt.Errorf("invalid parameter")
	}

	for i := 0; i < len(s); i++ {
		if s[i] != 0 {
			m.M[i] = m2.M[i]
		}
	}

	return nil
}

// 修改值为固定数值
func (m *Matrix) SetValueBySlice2(v float32, s []int) error {
	if len(s) > int(m.Row) {
		return fmt.Errorf("invalid parameter")
	}

	for i := 0; i < len(s); i++ {
		m.M[s[i]] = v
	}

	return nil
}

// 矩阵特定位置添加上所给值
func (m *Matrix) AddValueBySlice(v float32, s []int) error {
	if len(s) > int(m.Row) {
		return fmt.Errorf("invalid parameter")
	}
	for i := 0; i < len(s); i++ {
		m.M[s[i]] += v
	}
	return nil

}

// 返回第i列的数据
func (m *Matrix) GetColData(i int) []float32 {
	var result []float32
	for j := 0; j < int(m.Row); j++ {
		result = append(result, m.M[j*int(m.Col)+i])
	}
	return result
}

// 根据bool矩阵设置矩阵的值
// 坐标为true的地方设置为value
func (m *Matrix) SetValue(m1 []bool, value float32) {
	for i := 0; i < int(m.Row*m.Col); i++ {
		if m1[i] {
			m.M[i] = value
		}
	}
}

// 判断两个二维数组具体位置的大小情况
// m1 > m2
// 返回bool数组
func Bigger(m1, m2 Matrix) []bool {
	// 非法参数
	if !(m1.Row == m2.Row && m1.Col == m2.Col) {
		return []bool{}
	}

	result := make([]bool, m1.Row*m1.Col)

	for i := 0; i < int(m1.Row*m1.Col); i++ {

		if m1.M[i] > m2.M[i] {
			result[i] = true
		}

	}

	return result
}

// // 判断两个二维数组具体位置的大小情况
// // m1 < m2
// // 返回bool数组
// func Smaller(m1, m2 Matrix) [][]bool {
// 	// 非法参数
// 	if !(m1.Row == m2.Row && m1.Col == m2.Col) {
// 		return [][]bool{}
// 	}

// 	var result [][]bool

// 	for i := 0; i < int(m1.Row); i++ {
// 		temp := make([]bool, m1.Col)
// 		for j := 0; j < int(m1.Col); j++ {
// 			if m1.M[i][j] < m2.M[i][j] {
// 				temp[j] = true
// 			} else {
// 				temp[j] = false
// 			}
// 		}
// 		result = append(result, temp)
// 	}

// 	return result
// }

// 将灰度图片转换成灰度二维数组
func Convert(mat gocv.Mat) Matrix {
	var m Matrix
	m.Row = uint(mat.Rows())
	m.Col = uint(mat.Cols())

	// 图像的一维字节数组
	data := mat.ToBytes()

	// 转换为float32的一维数组
	var st []float32
	for j := 0; j < len(data); j++ {
		st = append(st, float32(data[j]))
	}

	m.M = st

	return m

}

// // 对数组进行切片操作
// func (m *Matrix) Slice(row uint) Matrix {
// 	var r Matrix
// 	r.Col = m.Col
// 	r.Row = row

// 	for i := 0; i < int(row); i++ {
// 		r.M = append(r.M, m.M[row])
// 	}

// 	return r

// }

// 将矩阵清零
func (m *Matrix) SetZero() {
	for i := 0; i < int(m.Row*m.Col); i++ {
		m.M[i] = 0
	}
}

// 计算矩阵所有元素的和
func Sum(m Matrix) float32 {
	var sum float32
	for i := 0; i < int(m.Row*m.Col); i++ {
		sum += m.M[i]
	}
	return sum

}

// 不修改传入矩阵，返回相加后结果
func Add(m1 Matrix, value float32) Matrix {
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)
	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] + value
	}
	matrix.M = M
	return matrix

}

// 不修改传入矩阵，返回相减后结果
func Sub(m1 Matrix, value float32) Matrix {
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)
	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] - value
	}
	matrix.M = M
	return matrix
}

// 不修改传入矩阵，返回相乘后结果
func Multiply(m1 Matrix, value float32) Matrix {
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)
	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] * value
	}
	matrix.M = M
	return matrix

}

// 不修改传入矩阵，返回相除后结果
func Divide(m1 Matrix, value float32) Matrix {
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)
	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] / value
	}
	matrix.M = M
	return matrix

}

// 矩阵相加，即对应位置元素相加
func AddMatrix(m1, m2 Matrix) (Matrix, error) {
	if (m1.Row != m2.Row) || (m1.Col != m2.Col) {
		return Matrix{}, fmt.Errorf("invalid parameter")
	}
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)

	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] + m2.M[i]
	}

	matrix.M = M

	return matrix, nil

}

// 矩阵相减，即对应位置元素相减
func SubMatrix(m1, m2 Matrix) (Matrix, error) {
	if (m1.Row != m2.Row) || (m1.Col != m2.Col) {
		return Matrix{}, fmt.Errorf("invalid parameter")
	}
	var matrix Matrix
	matrix.Row = m1.Row
	matrix.Col = m1.Col

	M := make([]float32, m1.Row*m1.Col)

	for i := 0; i < int(m1.Row*m1.Col); i++ {
		M[i] = m1.M[i] - m2.M[i]
	}

	matrix.M = M

	return matrix, nil

}

// 两个同行列的矩阵、对应位置的数字相乘
func MultiplyMatrix(m1 Matrix, m2 Matrix) (Matrix, error) {
	if m1.Row != m2.Row || m1.Col != m2.Col {
		return Matrix{}, fmt.Errorf("invalid matrix")
	}

	var matrix Matrix
	matrix.Col = m1.Col
	matrix.Row = m1.Row

	M := make([]float32, matrix.Row*matrix.Col)

	for i := 0; i < int(matrix.Row*matrix.Col); i++ {
		M[i] = m1.M[i] * m2.M[i]
	}
	matrix.M = M

	return matrix, nil

}

// 矩阵所有元素加上某值
func (m *Matrix) Add(data float32) {
	for i := 0; i < int(m.Row*m.Col); i++ {
		m.M[i] = m.M[i] + data
	}
}

// 矩阵所有元素减去某值
func (m *Matrix) Sub(data float32) {
	for i := 0; i < int(m.Row*m.Col); i++ {
		m.M[i] = m.M[i] - data
	}
}

// 矩阵所有元素除以某值
func (m *Matrix) Divide(data float32) {
	for i := 0; i < int(m.Row*m.Col); i++ {
		m.M[i] = m.M[i] / data
	}
}

// 所有元素乘上某值
func (m *Matrix) Multiply(data float32) {
	for i := 0; i < int(m.Row*m.Col); i++ {
		m.M[i] = m.M[i] * data
	}
}

// 返回一个[0, 1)的随机数
func RandomOne() float32 {
	return rand.Float32()
}

// 返回随机的二维数组
func Random(row uint, col uint) Matrix {

	var matrix Matrix
	matrix.Row = row
	matrix.Col = col

	M := make([]float32, row*col)

	for i := 0; i < int(row*col); i++ {
		M[i] = rand.Float32()
	}
	matrix.M = M
	return matrix
}

// 将矩阵以文件形式保存
func Save(path string, matrix Matrix) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(matrix)
	if err != nil {
		return err
	}
	data := buffer.Bytes()

	// 保存文件
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil

}

// 保存interface{}
func Save2(path string, obj interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(obj)
	if err != nil {
		return err
	}
	data := buffer.Bytes()

	// 保存文件
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil

}

// 加载保存的矩阵文件
func Load(path string) (Matrix, error) {
	file, err := os.Open(path)
	if err != nil {
		return Matrix{}, err
	}

	// 文件连接退出前关闭
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var matrix Matrix
	err = decoder.Decode(&matrix)
	if err != nil {
		return Matrix{}, err
	}

	return matrix, nil

}

// 不改变数据内容的情况下，改变一个数组的格式
func Reshape(m Matrix, row uint, col uint) (Matrix, error) {
	if m.Col*m.Row != row*col {
		return Matrix{}, fmt.Errorf("invalid parameter")
	}

	var matrix Matrix
	matrix.Row = row
	matrix.Col = col

	M := make([]float32, row*col)
	for i := 0; i < int(row*col); i++ {
		M[i] = m.M[i]
	}
	matrix.M = M

	return matrix, nil
}

// 返回bool切片
// 使用于二维矩阵
func WhereMatrix2(m Matrix, flag string, value float32) ([]bool, error) {
	if !(flag == ">" || flag == "<" || flag == ">=" || flag == "<=" || flag == "!=" || flag == "==") {
		return []bool{}, fmt.Errorf("invalid flag parameter")
	}

	var result []bool

	switch flag {
	case ">":
		// 返回满足(m > value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] > value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}
	case ">=":
		// 返回满足(m >= value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] >= value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}

	case "<":
		// 返回满足(m < value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] < value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}

	case "<=":
		// 返回满足(m <= value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] <= value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}

	case "!=":
		// 返回满足(m <= value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] != value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}

	case "==":
		// 返回满足(m == value)的下标
		for i := 0; i < int(m.Row*m.Col); i++ {

			if m.M[i] == value {
				result = append(result, true)
			} else {
				result = append(result, false)
			}

		}
	}

	return result, nil

}

// // 返回二维bool矩阵
// // 使用于二维矩阵
// func WhereMatrix(m Matrix, flag string, value float64) ([][]bool, error) {
// 	if !(flag == ">" || flag == "<" || flag == ">=" || flag == "<=" || flag == "!=" || flag == "==") {
// 		return [][]bool{}, fmt.Errorf("invalid flag parameter")
// 	}

// 	var result [][]bool

// 	switch flag {
// 	case ">":
// 		// 返回满足(m > value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] > value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}
// 	case ">=":
// 		// 返回满足(m >= value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] >= value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}

// 	case "<":
// 		// 返回满足(m < value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] < value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}

// 	case "<=":
// 		// 返回满足(m <= value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] <= value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}

// 	case "!=":
// 		// 返回满足(m <= value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] != value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}

// 	case "==":
// 		// 返回满足(m == value)的下标
// 		for i := 0; i < int(m.Row); i++ {
// 			temp := make([]bool, m.Col)
// 			for j := 0; j < int(m.Col); j++ {
// 				if m.M[i][j] == value {
// 					temp[j] = true
// 				} else {
// 					temp[j] = false
// 				}
// 			}
// 			result = append(result, temp)
// 		}
// 	}

// 	return result, nil

// }

// 只适用于列数为1的矩阵
// 返回满足给定条件的下标
func Where(m Matrix, flag string, value float32) ([]int, error) {
	// 只适用于列数为1的矩阵
	if m.Col != 1 {
		return []int{}, fmt.Errorf("invalid matrix parameter")
	}

	if !(flag == ">" || flag == "<" || flag == ">=" || flag == "<=" || flag == "!=") {
		return []int{}, fmt.Errorf("invalid flag parameter")
	}

	var result []int

	switch flag {
	case ">":
		// 返回满足(m > value)的下标
		for i := 0; i < int(m.Row); i++ {
			if m.M[i] > value {
				result = append(result, i)
			}
		}
	case ">=":
		// 返回满足(m >= value)的下标
		for i := 0; i < int(m.Row); i++ {
			if m.M[i] >= value {
				result = append(result, i)
			}
		}

	case "<":
		// 返回满足(m < value)的下标
		for i := 0; i < int(m.Row); i++ {
			if m.M[i] < value {
				result = append(result, i)
			}
		}

	case "<=":
		// 返回满足(m <= value)的下标
		for i := 0; i < int(m.Row); i++ {
			if m.M[i] <= value {
				result = append(result, i)
			}
		}

	case "!=":
		// 返回满足(m <= value)的下标
		for i := 0; i < int(m.Row); i++ {
			if m.M[i] != value {
				result = append(result, i)
			}
		}
	}

	return result, nil

}

// // 根据所给元素生成矩阵
// func Array(s ...[]float64) (Matrix, error) {
// 	// 矩阵行数
// 	row := len(s)

// 	if row == 0 {
// 		return Matrix{}, fmt.Errorf("nil input")
// 	}

// 	col := len(s[0])
// 	for i := 0; i < row; i++ {
// 		temp := len(s[i])
// 		if temp != col {
// 			return Matrix{}, fmt.Errorf("invalid input")
// 		}
// 	}

// 	matrix := Matrix{
// 		Row: uint(row),
// 		Col: uint(col),
// 		M:   s,
// 	}

// 	return matrix, nil

// }

// 生成全1矩阵
func Ones(row uint, col uint) Matrix {
	return newMatrix(row, col, 1)
}

// 生成全0矩阵
func Zeros(row uint, col uint) Matrix {
	return newMatrix(row, col, 0)
}

// 给定行、列, 生成数值全为data的矩阵
func newMatrix(row uint, col uint, data float32) Matrix {
	var matrix Matrix
	matrix.Row = row
	matrix.Col = col

	M := make([]float32, row*col)

	for i := 0; i < int(row*col); i++ {
		M[i] = data
	}
	matrix.M = M
	return matrix
}

// 返回float64切片里最大的值
func findMaxElement(arr []float64) float64 {
	max_num := arr[0]
	for i := 0; i < len(arr); i++ {
		if arr[i] > max_num {
			max_num = arr[i]
		}
	}
	return max_num
}
