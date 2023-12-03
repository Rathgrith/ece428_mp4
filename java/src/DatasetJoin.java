import java.io.IOException;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.input.MultipleInputs;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;

public class DatasetJoin {

    public static class JoinMapper1 extends Mapper<Object, Text, Text, Text> {
        private int joinColumnIndex;
        private Text outKey = new Text();
        private Text outValue = new Text();

        @Override
        protected void setup(Context context) throws IOException, InterruptedException {
            joinColumnIndex = context.getConfiguration().getInt("join.column.index1", 0);
        }

        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            String[] columns = value.toString().split(",");
            outKey.set(columns[joinColumnIndex]);
            outValue.set("set1," + value.toString());
            context.write(outKey, outValue);
        }
    }

    public static class JoinMapper2 extends Mapper<Object, Text, Text, Text> {
        private int joinColumnIndex;
        private Text outKey = new Text();
        private Text outValue = new Text();

        @Override
        protected void setup(Context context) throws IOException, InterruptedException {
            joinColumnIndex = context.getConfiguration().getInt("join.column.index2", 0);
        }

        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            String[] columns = value.toString().split(",");
            outKey.set(columns[joinColumnIndex]);
            outValue.set("set2," + value.toString());
            context.write(outKey, outValue);
        }
    }

    public static class JoinReducer extends Reducer<Text, Text, Text, Text> {
        private Text result = new Text();

        public void reduce(Text key, Iterable<Text> values, Context context) throws IOException, InterruptedException {
            // You can optimize this part for better performance and more complex join logic
            String set1Value = "";
            String set2Value = "";
            for (Text val : values) {
                String valueStr = val.toString();
                if (valueStr.startsWith("set1,")) {
                    set1Value = valueStr.substring(5);
                } else if (valueStr.startsWith("set2,")) {
                    set2Value = valueStr.substring(5);
                }
            }

            if (!set1Value.isEmpty() && !set2Value.isEmpty()) {
                result.set(set1Value + "," + set2Value);
                context.write(key, result);
            }
        }
    }

    public static void main(String[] args) throws Exception {
        if (args.length != 5) {
            System.err.println("Usage: DatasetJoin <input path1> <input path2> <output path> <join column index1> <join column index2>");
            System.exit(-1);
        }

        Configuration conf = new Configuration();
        conf.setInt("join.column.index1", Integer.parseInt(args[3]));
        conf.setInt("join.column.index2", Integer.parseInt(args[4]));

        Job job = Job.getInstance(conf, "Dataset Join");
        job.setJarByClass(DatasetJoin.class);
        job.setReducerClass(JoinReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class);

        MultipleInputs.addInputPath(job, new Path(args[0]), TextInputFormat.class, JoinMapper1.class);
        MultipleInputs.addInputPath(job, new Path(args[1]), TextInputFormat.class, JoinMapper2.class);
        FileOutputFormat.setOutputPath(job, new Path(args[2]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
