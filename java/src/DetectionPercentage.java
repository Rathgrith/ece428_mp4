import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

public class DetectionPercentage {

    public static class TokenizerMapper extends Mapper<Object, Text, Text, IntWritable> {

        private final static IntWritable one = new IntWritable(1);
        private Text word = new Text();
        private String interconneType;

        @Override
        protected void setup(Context context) throws IOException, InterruptedException {
            interconneType = context.getConfiguration().get("interconneType");
        }

        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            String[] dataArray = value.toString().split(","); // split the data into array
            if (dataArray.length > 10) { // avoid null pointer exception
                if (dataArray[10].trim().equals(interconneType)) { // check interconne type at index 10
                    word.set(dataArray[9]); // set 'Detection_' value from index 9
                    context.write(word, one);
                }
            }
        }
    }

    public static class IntSumReducer extends Reducer<Text, IntWritable, Text, Text> {

        private Map<String, Integer> countMap = new HashMap<>();

        public void reduce(Text key, Iterable<IntWritable> values, Context context) throws IOException, InterruptedException {
            int sum = 0;
            for (IntWritable val : values) {
                sum += val.get();
            }
            countMap.put(key.toString(), sum);
        }

        @Override
        protected void cleanup(Context context) throws IOException, InterruptedException {
            int totalCount = countMap.values().stream().mapToInt(Integer::intValue).sum();
            for (Map.Entry<String, Integer> entry : countMap.entrySet()) {
                double percentage = 100.0 * entry.getValue() / totalCount;
                context.write(new Text(entry.getKey()), new Text(percentage + "%"));
            }
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        conf.set("interconneType", args[2]);
        Job job = Job.getInstance(conf, "detection percentage");
        job.setJarByClass(DetectionPercentage.class);
        job.setMapperClass(TokenizerMapper.class);
        job.setCombinerClass(IntSumReducer.class);
        job.setReducerClass(IntSumReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class); // Set output value class to Text
        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
